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
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gizak/termui"
)

func RenderBucketExplorerListing(bucket BucketWithDisplay, nodes []*Node, selection int, deferFunc func()) {

	// Get a UI ready list depending on where the selection pointer is

	list := CreateDirectoryList(bucket.displayString, nodes, selection)
	termui.Clear()
	termui.Render(list, RenderHelp())

	// Set default handlers and defer tempdir removal

	SetDefaultHandlers(deferFunc)

	// Up key moves up

	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		if selection == 0 {
			return
		} else {
			selection -= 1
			list := CreateDirectoryList(bucket.displayString, nodes, selection)
			termui.Clear()
			termui.Render(list, RenderHelp())
		}
	})

	// Down key moves down

	termui.Handle("/sys/kbd/<down>", func(termui.Event) {
		if selection == len(nodes)-1 {
			return
		} else {
			selection += 1
			list := CreateDirectoryList(bucket.displayString, nodes, selection)
			termui.Clear()
			termui.Render(list, RenderHelp())
		}
	})

	// Enter descends the directory or initiates a download for a file

	termui.Handle("/sys/kbd/<enter>", func(termui.Event) {

		if nodes[selection].Info.IsDir {

			// A directory was selected

			log.Printf("Selected Directory: %s\n", nodes[selection].DisplayString)

			SetBackHandler(func() {
				termui.ResetHandlers()
				if nodes[0].Parent == nil {
					log.Println("Reached directory root, returning to buckets")
					deferFunc()
					ReloadMainBuckets()
				} else {
					log.Printf("Going back to directory: %+v\n", nodes[0].Parent.DisplayString)
					listing := GetNodeDirectory(nodes[0].Parent)
					RenderBucketExplorerListing(bucket, listing, 0, deferFunc)
				}
			})

			// Call parent function for selected node

			log.Printf("Descending into node: %+v\n", nodes[selection].DisplayString)
			listing := GetNodeDirectory(nodes[selection])
			RenderBucketExplorerListing(bucket, listing, 0, deferFunc)

		} else {

			// A File was selected

			log.Printf("File Selected: %s\n", nodes[selection].DisplayString)

			// Current downloads default to working directory

			dest := filepath.Join(currentWorkingDir, path.Base(*nodes[selection].S3Object.Key))
			p := CreateDownloadPrompt(dest)
			termui.Render(p)

			// Create an AWS Session in the region of the bucket

			sess, err := InitSession(bucket.region)
			if err != nil {
				log.Println(err)
				RenderError(err.Error())
				return
			}

			// Download the file

			err = sess.DownloadObject(bucket, nodes[selection], dest)
			if err != nil {
				log.Println(err)
				RenderError(err.Error())
			} else {
				p := CreateFinishedDownloadPrompt(dest)
				termui.Render(p)
			}
		}

	})

}

func RenderBucketExplorer(bucket BucketWithDisplay) {

	// init selection pointer

	var selection int
	selection = 0

	// retrieve all objects for bucket

	objects, err := s3Session.GetBucketObjects(bucket)
	if err != nil {
		ReloadMainBucketsWithError(err)
		return

	}

	// create a local mock filesystem for easier indexing

	mockFsRoot, err := CreateMockFs(objects)
	if err != nil {
		ReloadMainBucketsWithError(err)
	}

	// Set clean function to pass to all default handlers

	deferFunc := func() {
		log.Printf("Cleaning Temp Directory: %s\n", mockFsRoot)
		err := os.RemoveAll(mockFsRoot)
		if err != nil {
			log.Printf("Error cleaning Temp (%s): %s\n", mockFsRoot, err)
		}
	}

	// Evaluate the directory tree

	tree, err := NewTree(objects, mockFsRoot)
	if err != nil {
		ReloadMainBucketsWithError(err)
	}

	// Get a listing for the root node

	listing := GetNodeDirectory(tree)
	if err != nil {
		ReloadMainBucketsWithError(err)
	}

	// Render the bucket explorer

	SetBackHandler(ReloadMainBuckets, deferFunc)
	RenderBucketExplorerListing(bucket, listing, selection, deferFunc)

}
