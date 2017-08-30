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
	"fmt"
	"os"
	"path/filepath"
)

//Build performs (if needed) the build of the given package into the given
//target directory.
func Build(pkg Package, cache Cache, targetDirPath string) error {
	entry, err := cache.GetEntryFor(pkg)
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
