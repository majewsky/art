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
)

func main() {
	os.Exit(_main())
}

func _main() (exitCode int) {
	ui := &UI{}

	cfg, err := readConfig()
	if err != nil {
		ui.ShowError(err)
		return 1
	}
	mcfg, err := readMakepkgConfig()
	if err != nil {
		ui.ShowError(err)
		return 1
	}

	cache, err := readCache()
	if err != nil {
		ui.ShowError(err)
		return 1
	}

	ui.SetCurrentTask("Discovering packages", uint(len(cfg.Sources)))
	var totalPackageCount uint
	for _, src := range cfg.Sources {
		err := src.discoverPackages(mcfg)
		if err != nil {
			ui.ShowError(err)
			exitCode = 1
		}
		ui.StepTask()
		totalPackageCount += uint(len(src.Packages))
	}
	ui.EndTask()

	if exitCode > 0 {
		return
	}

	ui.SetCurrentTask("Building packages", totalPackageCount)
	for _, src := range cfg.Sources {
		for _, pkg := range src.Packages {
			err := cache.Build(pkg, cfg.Target.Path, ui)
			if err != nil {
				ui.ShowError(err)
				exitCode = 1
			}
			ui.StepTask()
		}
	}
	ui.EndTask()

	err = cache.writeCache()
	if err != nil {
		ui.ShowError(err)
		exitCode = 1
	}

	if exitCode > 0 {
		return
	}

	ui.SetCurrentTask("Post-processing and signing packages", totalPackageCount)
	var allOutputFiles []string
	for _, src := range cfg.Sources {
		for _, pkg := range src.Packages {
			files, err := cache.AddMissingSignatures(pkg, cfg.Target.Path, mcfg)
			if err != nil {
				ui.ShowError(err)
				exitCode = 1
			}
			allOutputFiles = append(allOutputFiles, files...)
			ui.StepTask()
		}
	}
	ui.EndTask()
	if exitCode > 0 {
		return
	}

	ok := cfg.Target.addNewPackages(allOutputFiles, cache, ui)
	if !ok {
		exitCode = 1
		return
	}
	err = cache.writeCache() //since the previous call might have changed it
	if err != nil {
		ui.ShowError(err)
		exitCode = 1
	}
	ok = cfg.Target.pruneMetadata(allOutputFiles, ui)
	if !ok {
		exitCode = 1
		return
	}
	ok = cfg.Target.prunePackages(allOutputFiles, ui)
	if !ok {
		exitCode = 1
		return
	}

	return
}
