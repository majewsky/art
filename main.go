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
)

func main() {
	os.Exit(_main())
}

func _main() (exitCode int) {
	cfg, err := readConfig()
	if err != nil {
		showError(err)
		return 1
	}
	mcfg, err := readMakepkgConfig()
	if err != nil {
		showError(err)
		return 1
	}
	fmt.Printf("mcfg = %#v\n", mcfg)

	cache, err := readCache()
	if err != nil {
		showError(err)
		return 1
	}

	progress("Discovering packages")
	for _, src := range cfg.Sources {
		err := src.discoverPackages(mcfg)
		if err != nil {
			showError(err)
			exitCode = 1
		}
		step()
	}

	if exitCode > 0 {
		return
	}

	progress("Building packages")
	for _, src := range cfg.Sources {
		for _, pkg := range src.Packages {
			err := cache.Build(pkg, cfg.Target.Path)
			if err != nil {
				showError(err)
				exitCode = 1
			}
			step()
		}
	}
	done()

	err = cache.writeCache()
	if err != nil {
		showError(err)
		exitCode = 1
	}

	if exitCode > 0 {
		return
	}

	return
}

func showError(err error) {
	if err != nil {
		fmt.Printf("\x1B[1;31m!! \x1B[0;31m%s\x1B[0m\n", err.Error())
	}
}

func progress(msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	fmt.Printf("\x1B[1;36m>> \x1B[0;36m%s\x1B[0m", msg)
}

func step() {
	fmt.Printf("\x1B[0;36m.\x1B[0m")
}

func done() {
	fmt.Printf("\n")
}
