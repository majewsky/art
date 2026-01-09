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
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

// MakepkgConfig contains the fields from makepkg.conf that interest us.
type MakepkgConfig struct {
	Architecture string
	GPGKeyID     string
}

var makepkgFieldRegex = regexp.MustCompile(`^([A-Z]+)\s*=\s*["']?(.*?)["']?$`)

func readMakepkgConfig() (MakepkgConfig, error) {
	bytes, err := ioutil.ReadFile("/etc/makepkg.conf")
	if err != nil {
		return MakepkgConfig{}, err
	}

	var result MakepkgConfig
	result.GPGKeyID = os.Getenv("GPGKEY")
	for _, line := range strings.Split(string(bytes), "\n") {
		line = strings.TrimSpace(line)
		match := makepkgFieldRegex.FindStringSubmatch(line)
		if match != nil {
			switch match[1] {
			case "CARCH":
				result.Architecture = match[2]
			case "GPGKEY":
				result.GPGKeyID = match[2]
			}
		}
	}

	return result, nil
}

// FilterFilesForCurrentArch takes a list of output files that can be generated
// by a PKGBUILD, and returns only these matching the current architecture (i.e.
// where the architecture is the current one or "any").
func (cfg MakepkgConfig) FilterFilesForCurrentArch(outputFiles []string) (result []string) {
	arch1 := "-" + cfg.Architecture
	arch2 := "-any"

	for _, outputFile := range outputFiles {
		str := strings.TrimSuffix(outputFile, ".pkg.tar.xz")
		if strings.HasSuffix(outputFile, ".zst") {
			str = strings.TrimSuffix(outputFile, ".pkg.tar.zst")
		}
		if strings.HasSuffix(str, arch1) || strings.HasSuffix(str, arch2) {
			result = append(result, outputFile)
		}
	}
	return
}
