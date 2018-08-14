/**
This file is part of s3explorer.

s3explorer is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

s3explorer is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with s3explorer.  If not, see <https://www.gnu.org/licenses/>.
**/

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gizak/termui"
)

func SetDefaultHandlers(exitFunc func()) {

	// Always set "q" to exit with any relevant deferred functions

	SetExitHandler(exitFunc)
}

func SetExitHandler(exitFunc func()) {

	// Run a prepared function and exit the program

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		log.Println("Received User Requested Exit")
		log.Println("Running extra handlers")
		exitFunc()
		log.Println("Stopping Event loop")
		termui.StopLoop()
		fmt.Println()
		os.Exit(EXIT_USER_REQUESTED)
	})
}

func SetBackHandler(runFuncs ...func()) {

	// Set the back handler to a prepared function

	termui.Handle("/sys/kbd/b", func(termui.Event) {
		log.Println("Received <back> request, running handlers")
		for _, runFunc := range runFuncs {
			runFunc()
		}
	})
}
