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
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Package is a package definition. It is satisfied by types HoloBuildPackage
// and NativePackage.
type Package interface {
	//CacheKey returns a string that uniquely idenfities this package.
	CacheKey() string
	//LastModified returns the mtime of the package definition file.
	LastModified() (time.Time, error)
	//OutputFiles returns the list of files produced by Build().
	OutputFiles() ([]string, error)
	//Build builds all output files into the given target directory.
	Build(targetDirPath string) error
}

// HoloBuildPackage describes a package declaration that can be built by using
// holo-build(8).
type HoloBuildPackage struct {
	Path          string
	MakepkgConfig MakepkgConfig
}

// CacheKey implements the Package interface.
func (pkg HoloBuildPackage) CacheKey() string {
	return pkg.Path
}

// LastModified implements the Package interface.
func (pkg HoloBuildPackage) LastModified() (time.Time, error) {
	fi, err := os.Stat(pkg.Path)
	if err != nil {
		return time.Time{}, err
	}
	return fi.ModTime(), nil
}

// OutputFiles implements the Package interface.
func (pkg HoloBuildPackage) OutputFiles() ([]string, error) {
	cmd := exec.Command("holo-build", "--suggest-filename", pkg.Path)
	cmd.Stdin = nil
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	result := []string{strings.TrimSpace(string(buf.Bytes()))}
	return pkg.MakepkgConfig.FilterFilesForCurrentArch(result), err
}

// Build implements the Package interface.
func (pkg HoloBuildPackage) Build(targetDirPath string) error {
	absPath, err := filepath.Abs(pkg.Path)
	if err != nil {
		return err
	}

	cmd := exec.Command("holo-build", absPath)
	cmd.Dir = targetDirPath
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// NativePackage describes a directory with a PKGBUILD that can be built using
// makepkg(8).
type NativePackage struct {
	Path          string
	MakepkgConfig MakepkgConfig
}

// CacheKey implements the Package interface.
func (pkg NativePackage) CacheKey() string {
	return pkg.Path
}

// LastModified implements the Package interface.
func (pkg NativePackage) LastModified() (time.Time, error) {
	fi, err := os.Stat(pkg.Path)
	if err != nil {
		return time.Time{}, err
	}
	return fi.ModTime(), nil
}

// OutputFiles implements the Package interface.
func (pkg NativePackage) OutputFiles() ([]string, error) {
	cmd := exec.Command("makepkg", "--packagelist", "-p", filepath.Base(pkg.Path))
	cmd.Dir = filepath.Dir(pkg.Path)
	cmd.Stdin = nil
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	var result []string
	for _, line := range strings.Split(string(buf.Bytes()), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		result = append(result, filepath.Base(line))
	}
	return pkg.MakepkgConfig.FilterFilesForCurrentArch(result), nil
}

// Build implements the Package interface.
func (pkg NativePackage) Build(targetDirPath string) error {
	cmd := exec.Command("makepkg", "-s", "-p", filepath.Base(pkg.Path))
	cmd.Dir = filepath.Dir(pkg.Path)
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"PKGDEST="+targetDirPath,
	)
	return cmd.Run()
}
