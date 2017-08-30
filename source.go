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
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//Source describes a directory from which packages are read.
type Source struct {
	Path     string    `toml:"path"`
	Packages []Package `toml:"-"`
}

func (s *Source) discoverPackages(mcfg MakepkgConfig) error {
	dir, err := os.Open(s.Path)
	if err != nil {
		return err
	}
	fis, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	sort.Slice(fis, func(i, j int) bool { return fis[i].Name() < fis[j].Name() })

	for _, fi := range fis {
		pkgPath := filepath.Join(s.Path, fi.Name())

		//subdirectory: is a package if there is a PKGBUILD
		if fi.IsDir() {
			fi2, err := os.Stat(filepath.Join(pkgPath, "PKGBUILD"))
			if err != nil {
				if os.IsNotExist(err) {
					continue
				}
				return err
			}
			if isRegularOrSymlink(fi2.Mode()) {
				s.Packages = append(s.Packages, &NativePackage{
					Path:          pkgPath,
					MakepkgConfig: mcfg,
				})
			}
		}

		//regular file or symlink: is a holo-build package if suffix ".pkg.toml"
		if isRegularOrSymlink(fi.Mode()) && strings.HasSuffix(pkgPath, ".pkg.toml") {
			s.Packages = append(s.Packages, &HoloBuildPackage{
				Path:          pkgPath,
				MakepkgConfig: mcfg,
			})
		}
	}

	return nil
}

func isRegularOrSymlink(mode os.FileMode) bool {
	return mode.IsRegular() || (mode&os.ModeType) == os.ModeSymlink
}
