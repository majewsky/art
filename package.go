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
)

//Package is a package definition. It is satisfied by types HoloBuildPackage
//and NativePackage.
type Package interface {
	//Build builds this package, stores the resulting package files in
	//`targetDirPath` and returns a list of their file names.
	Build(targetDirPath string) ([]string, error)
}

//HoloBuildPackage describes a package declaration that can be built by using
//holo-build(8).
type HoloBuildPackage struct {
	Path string
}

//Build implements the Package interface.
func (pkg HoloBuildPackage) Build(targetDirPath string) ([]string, error) {
	absPath, err := filepath.Abs(pkg.Path)
	if err != nil {
		return nil, err
	}

	//NOTE: It would be nice if `holo-build` could write the package to the
	//target directory, and print its name to stdout in a single pass.
	cmd := exec.Command("holo-build", "--suggest-filename", absPath)
	cmd.Dir = targetDirPath
	cmd.Stdin = nil
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}
	outputFileName := strings.TrimSpace(string(buf.Bytes()))

	cmd = exec.Command("holo-build", absPath)
	cmd.Dir = targetDirPath
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return []string{outputFileName}, cmd.Run()
}

//NativePackage describes a directory with a PKGBUILD that can be built using
//makepkg(8).
type NativePackage struct {
	Path string
}

//Build implements the Package interface.
func (pkg NativePackage) Build(targetDirPath string) ([]string, error) {
	//TODO: unimplemented
	return nil, nil
}
