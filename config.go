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
	"errors"
	"fmt"
	"io/ioutil"

	toml "github.com/BurntSushi/toml"
)

//Configuration is the contents of the configuration file.
type Configuration struct {
	Sources []*Source  `toml:"source"`
	Target  Repository `toml:"target"`
}

func readConfig() (*Configuration, error) {
	bytes, err := ioutil.ReadFile("./art.toml")
	if err != nil {
		return nil, err
	}

	var cfg Configuration
	err = toml.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Target.Path == "" {
		return nil, errors.New("parse art.toml: missing value for target.path")
	}
	if cfg.Target.Name == "" {
		return nil, errors.New("parse art.toml: missing value for target.name")
	}
	if len(cfg.Sources) == 0 {
		return nil, errors.New("parse art.toml: no sources specified")
	}
	for idx, src := range cfg.Sources {
		if src.Path == "" {
			return nil, fmt.Errorf("parse art.toml: missing value for sources[%d].path", idx)
		}
	}

	return &cfg, nil
}
