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
	"crypto/md5"
	"encoding/hex"
	"os"
	"time"
)

// Returns true if the two times differ by less than a second.
func fuzzyTimeEqual(t1 time.Time, t2 time.Time) bool {
	unix1 := t1.UnixNano()
	unix2 := t2.UnixNano()
	diff := unix1 - unix2
	return diff > -1e9 && diff < 1e9
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	switch {
	case err == nil:
		return true, nil
	case os.IsNotExist(err):
		return false, nil
	default:
		return false, err
	}
}

func md5digest(buf []byte) string {
	s := md5.Sum(buf)
	return hex.EncodeToString(s[:])
}
