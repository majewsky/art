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

import "fmt"

//UI encapsulates the state of the terminal display.
type UI struct{}

//ShowError prints the given error if it is not nil.
func (ui *UI) ShowError(err error) {
	if err != nil {
		fmt.Printf("\x1B[1;31m!! \x1B[0;31m%s\x1B[0m\n", err.Error())
	}
}

//SetCurrentTask displays the progress of the next task.
func (ui *UI) SetCurrentTask(task string, count uint) {
	fmt.Printf("\x1B[1;36m>> \x1B[0;36m%s\x1B[0m", task)
}

//StepTask increases the counter on the task.
func (ui *UI) StepTask() {
	fmt.Printf("\x1B[0;36m.\x1B[0m")
}

//EndTask signals the end of the current task.
func (ui *UI) EndTask() {
	fmt.Printf("\n")
}
