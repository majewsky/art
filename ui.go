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

//if it looks stupid and works, it ain't stupid
const clearLine = "\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08\x08"

//UI encapsulates the state of the terminal display.
type UI struct {
	task  string
	step  uint
	count uint
}

//ShowError prints the given error if it is not nil.
func (ui *UI) ShowError(err error) {
	if err != nil {
		if ui.task != "" {
			fmt.Printf("\n")
		}
		fmt.Printf("\x1B[1;31m[error] \x1B[0;31m%s\x1B[0m\n", err.Error())
	}
}

//ShowWarning prints the given warning.
func (ui *UI) ShowWarning(msg string, args ...interface{}) {
	if ui.task != "" {
		fmt.Printf("\n")
	}
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	fmt.Printf("\x1B[1;33m[ warn] \x1B[0;33m%s\x1B[0m\n", msg)
}

//SetCurrentTask displays the progress of the next task.
func (ui *UI) SetCurrentTask(task string, count uint) {
	if ui.task != "" {
		ui.EndTask()
	}
	ui.task = task
	ui.step = 0
	ui.count = count
	ui.displayTask()
}

//StepTask increases the counter on the task.
func (ui *UI) StepTask() {
	ui.step++
	ui.displayTask()
}

//EndTask signals the end of the current task.
func (ui *UI) EndTask() {
	if ui.task != "" {
		ui.step = ui.count
		ui.displayTask()
		fmt.Printf("\n")

		ui.task = ""
		ui.step = 0
		ui.count = 0
	}
}

func (ui *UI) displayTask() {
	fmt.Printf(clearLine)
	progress := "....."
	if ui.count > 0 {
		progress = fmt.Sprintf("%2d/%2d", ui.step, ui.count)
	}
	fmt.Printf("\x1B[1;36m[%s] \x1B[0;36m%s\x1B[0m ", progress, ui.task)
}
