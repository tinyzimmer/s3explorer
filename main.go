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
)

const (
	VERSION = "v0.1"

	// Dump logs to /dev/null by default
	DEFAULT_LOG_FILE = os.DevNull

	// FS Options
	DEFAULT_DIRECTORY_MODE = 0755
	DEFAULT_FILE_MODE      = 0644

	// Exit Codes
	EXIT_USER_REQUESTED        = 0 // user exit
	EXIT_FAILED_NO_TERMINAL    = 1 // termui.Init() failed
	EXIT_FAILED_NO_LOGGER      = 2 // Unable to access log file or /dev/null
	EXIT_FAILED_AWS_CONNECT    = 3 // Could not connect to AWS S3 API
	EXIT_FAILED_BUCKET_LISTING = 4 // Could not get initial bucket listing

	// UI Options
	RIGHT_BUFFER              = 10
	LOWER_BUFFER              = 10
	CHECK_TERM_SLEEP_INTERVAL = 1
	MIN_TERM_HEIGHT_REQUIRED  = 15

	// AWS Options
	DEFAULT_REGION = "us-west-2" // Used for root-level ListBuckets operations
)

var (
	s3Session         S3Session // initial s3 session
	localDelimiter    string    // local filesystem path delimiter
	logFile           string    // log file
	currentWorkingDir string    // starting local working directory
	versionDump       bool      // version dump
)

func dumpVersion() {

	// Print the version and exit

	fmt.Printf("s3explorer version: %s\n", VERSION)
	os.Exit(EXIT_USER_REQUESTED)
}

func init() {

	// Debug will print a chatty logfile

	flag.StringVar(&logFile, "d", DEFAULT_LOG_FILE, "Path to write debug logs")
	flag.BoolVar(&versionDump, "v", false, "Print version and exit")
	flag.Parse()

	if versionDump {
		dumpVersion()
	}

	var err error
	var logWriter io.Writer

	if FileExists(logFile) && logFile != os.DevNull {

		// Remove pre-existing log file if exists

		os.Remove(logFile)
		logWriter, err = os.Create(logFile)
		if err != nil {
			fmt.Printf("Error: Failed to create log file %s\n", logFile)
			os.Exit(EXIT_FAILED_NO_LOGGER)
		}

	} else {

		// Otherwise just open the log file

		logWriter, err = os.Create(logFile)
		if err != nil {
			fmt.Printf("Error: Failed to create log file %s\n", logFile)
			os.Exit(EXIT_FAILED_NO_LOGGER)
		}

	}

	// Set the logger
	log.SetOutput(logWriter)
	log.Println("Started Debug Log")

	// Get Current Working Directory if we can (we start here for saving files)

	currentWorkingDir, err = os.Getwd()
	if err != nil {
		log.Println("Error: (non-fatal) Could not get working directory")
	} else {
		log.Printf("Got current working directory: %s\n", currentWorkingDir)
	}

	// Create an initial s3 session for bucket listing
	//		ListBuckets returns buckets for all regions

	s3Session, err = InitSession(DEFAULT_REGION)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(EXIT_FAILED_AWS_CONNECT)
	}

	// Get the local delimiter.
	// It's actually safe to use a POSIX path delimiter on Windows, but this feels safer

	localDelimiter = GetLocalDelimiter()
	log.Println("Finished init")
}

func main() {

	// Start the UI

	RunUi()

}
