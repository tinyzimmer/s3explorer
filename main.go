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
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gizak/termui"
)

const (
	VERSION = "v0.1"

	// Dump logs to /dev/null by default
	DEFAULT_LOG_FILE = os.DevNull

	// FS Options
	DEFAULT_DIRECTORY_MODE = 0755
	DEFAULT_FILE_MODE      = 0644

	// Exit Codes
	EXIT_USER_REQUESTED        = 0
	EXIT_FAILED_NO_TERMINAL    = 1
	EXIT_FAILED_NO_LOGGER      = 2
	EXIT_FAILED_AWS_CONNECT    = 3
	EXIT_FAILED_BUCKET_LISTING = 4

	// UI Options
	RIGHT_BUFFER = 10

	// AWS Options
	DEFAULT_REGION = "us-west-2"
)

var (
	s3Session         S3Session
	localDelimiter    string
	logFile           string
	currentWorkingDir string
	versionDump       bool
)

func dumpVersion() {
	fmt.Printf("s3explorer version: %s\n", VERSION)
	os.Exit(EXIT_USER_REQUESTED)
}

func init() {

	flag.StringVar(&logFile, "d", DEFAULT_LOG_FILE, "Path to write debug logs")
	flag.BoolVar(&versionDump, "v", false, "Print version and exit")
	flag.Parse()

	if versionDump {
		dumpVersion()
	}

	var err error
	var logWriter io.Writer

	if FileExists(logFile) {
		os.Remove(logFile)
		logWriter, err = os.Create(logFile)
		if err != nil {
			fmt.Printf("Error: Failed to create log file %s\n", logFile)
			os.Exit(EXIT_FAILED_NO_LOGGER)
		}
	} else {
		logWriter, err = os.Create(logFile)
		if err != nil {
			fmt.Printf("Error: Failed to create log file %s\n", logFile)
			os.Exit(EXIT_FAILED_NO_LOGGER)
		}
	}
	log.SetOutput(logWriter)
	log.Println("Started Debug Log")

	currentWorkingDir, err = os.Getwd()
	if err != nil {
		log.Println("Error: (non-fatal) Could not get working directory")
	} else {
		log.Printf("Got current working directory: %s\n", currentWorkingDir)
	}

	s3Session, err = InitSession(DEFAULT_REGION)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(EXIT_FAILED_AWS_CONNECT)
	}
	localDelimiter = GetLocalDelimiter()
	log.Println("Finished init")
}

func main() {
	RunUi()
}

func RunUi() {
	err := termui.Init()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(EXIT_FAILED_NO_TERMINAL)
	}
	defer termui.Close()
	buckets, err := s3Session.GetBucketWithDisplayStrings()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(EXIT_FAILED_BUCKET_LISTING)
	}
	SetDefaultHandlers(func() { return })
	RenderBucketListing(buckets)
	termui.Loop()
}
