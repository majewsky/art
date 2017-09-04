/*******************************************************************************
*
* Copyright 2017 Stefan Majewsky <majewsky@gmx.net>
*
* This program is free software: you can redistribute it and/or modify it under
* the terms of the GNU General Public License as published by the Free Software
* Foundation, either version 3 of the License, or (at your option) any later
* version.
*
* This program is distributed in the hope that it will be useful, but WITHOUT ANY
* WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
* A PARTICULAR PURPOSE. See the GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License along with
* this program. If not, see <http://www.gnu.org/licenses/>.
*
*******************************************************************************/

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

//PackageCacheEntry contains metadata for a Package instance that is held in the Cache.
type PackageCacheEntry struct {
	LastModified time.Time
	OutputFiles  []string
}

//OutputCacheEntry contains metadata for an output file that is held in the Cache.
type OutputCacheEntry struct {
	MD5Digest string
}

//Cache contains metadata for a number of Package instances.
type Cache struct {
	Packages    map[string]PackageCacheEntry `toml:"package"`
	OutputFiles map[string]OutputCacheEntry  `toml:"output"`
}

const (
	cachePath = ".art-cache"
)

func readCache() (*Cache, error) {
	c := &Cache{
		Packages:    make(map[string]PackageCacheEntry),
		OutputFiles: make(map[string]OutputCacheEntry),
	}

	bytes, err := ioutil.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			//acceptable, e.g. on first run; start with empty cache
			return c, nil
		}
		return nil, err
	}

	err = toml.Unmarshal(bytes, c)
	return c, err
}

func (c *Cache) writeCache() error {
	var buf bytes.Buffer
	err := toml.NewEncoder(&buf).Encode(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cachePath, buf.Bytes(), 0644)
}

//GetEntryForPackage retrieves (or creates) a cache entry for the given Package.
func (c *Cache) GetEntryForPackage(pkg Package) (PackageCacheEntry, error) {
	entry, exists := c.Packages[pkg.CacheKey()]

	mtime, err := pkg.LastModified()
	if err != nil {
		return PackageCacheEntry{}, err
	}
	if exists && fuzzyTimeEqual(entry.LastModified, mtime) {
		return entry, nil
	}

	entry = PackageCacheEntry{
		LastModified: mtime,
	}
	entry.OutputFiles, err = pkg.OutputFiles()
	if err != nil {
		return PackageCacheEntry{}, err
	}

	c.Packages[pkg.CacheKey()] = entry
	return entry, nil
}

//GetEntryForOutputFile retrieves (or creates) a cache entry for the given output file.
func (c *Cache) GetEntryForOutputFile(path string) (OutputCacheEntry, error) {
	baseName := filepath.Base(path)
	entry, exists := c.OutputFiles[baseName]
	if exists {
		return entry, nil
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return OutputCacheEntry{}, err
	}

	entry = OutputCacheEntry{
		MD5Digest: md5digest(buf),
	}
	c.OutputFiles[baseName] = entry
	return entry, nil
}

////////////////////////////////////////////////////////////////////////////////

//Build performs (if needed) the build of the given package into the given
//target directory.
func (c *Cache) Build(pkg Package, targetDirPath string) error {
	entry, err := c.GetEntryForPackage(pkg)
	if err != nil {
		return err
	}

	var (
		alreadyBuilt = false
		needsBuild   = false
	)

	for _, fileName := range entry.OutputFiles {
		fi, err := os.Stat(filepath.Join(targetDirPath, fileName))
		switch {
		case err == nil:
			if fi.ModTime().After(entry.LastModified) {
				alreadyBuilt = true
			} else {
				return fmt.Errorf(
					"refusing to build %s: target file exists and is older than package definition",
					fileName,
				)
			}
		case os.IsNotExist(err):
			needsBuild = true
		default:
			return err
		}
	}

	if alreadyBuilt && needsBuild {
		return fmt.Errorf(
			"cannot build package: some of %v exist at target, but some do not",
			entry.OutputFiles,
		)
	}

	if !needsBuild {
		return nil
	}

	return pkg.Build(targetDirPath)
}

//AddMissingSignatures adds signature files to all output files that do not
//have one yet. It returns a list of the names of all output files.
func (c *Cache) AddMissingSignatures(pkg Package, targetDirPath string, mcfg MakepkgConfig) ([]string, error) {
	entry, err := c.GetEntryForPackage(pkg)
	if err != nil {
		return nil, err
	}

	if mcfg.GPGKeyID != "" {
		for _, fileName := range entry.OutputFiles {
			path := filepath.Join(targetDirPath, fileName)
			outputExists, err := fileExists(path)
			if err != nil {
				return nil, err
			}
			if !outputExists {
				continue
			}

			signatureExists, err := fileExists(path + ".sig")
			if err != nil {
				return nil, err
			}
			if signatureExists {
				continue
			}
			cmd := exec.Command(
				"gpg", "--detach-sign", "--use-agent",
				"-u", mcfg.GPGKeyID,
				"--no-armor", path,
			)
			cmd.Stdin = nil
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				return nil, err
			}
		}
	}

	return entry.OutputFiles, nil
}
