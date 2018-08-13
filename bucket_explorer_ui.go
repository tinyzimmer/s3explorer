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

	"github.com/gizak/termui"
)

func RenderBucketExplorerListing(bucket BucketWithDisplay, nodes []*Node, selection int, deferFunc func()) {

	list := CreateDirectoryList(bucket.displayString, nodes, selection)
	termui.Clear()
	termui.Render(list, RenderHelp())

	SetDefaultHandlers(deferFunc)

	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		if selection == 0 {
			return
		} else {
			selection -= 1
			list := CreateDirectoryList(bucket.displayString, nodes, selection)
			termui.Render(list, RenderHelp())
		}
	})

	termui.Handle("/sys/kbd/<down>", func(termui.Event) {
		if selection == len(nodes)-1 {
			return
		} else {
			selection += 1
			list := CreateDirectoryList(bucket.displayString, nodes, selection)
			termui.Render(list, RenderHelp())
		}
	})

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
					log.Printf("Going back to directory: %+v\n", nodes[0].Parent)
					listing := GetNodeDirectory(nodes[0].Parent)
					RenderBucketExplorerListing(bucket, listing, 0, deferFunc)
				}
			})

			log.Printf("Descending into node: %+v\n", nodes[selection])
			listing := GetNodeDirectory(nodes[selection])
			RenderBucketExplorerListing(bucket, listing, 0, deferFunc)

		} else {

			// A File was selected

			log.Printf("File Selected: %s\n", nodes[selection].DisplayString)

			p := CreateDownloadPrompt(nodes[selection])
			termui.Render(p)

		}

	})

}

func RenderBucketExplorer(bucket BucketWithDisplay) {

	var selection int

	objects, err := s3Session.GetBucketObjects(bucket)
	if err != nil {
		ReloadMainBucketsWithError(err)
		return

	}

	mockFsRoot, err := CreateMockFs(objects)
	if err != nil {
		ReloadMainBucketsWithError(err)
	}

	deferFunc := func() {
		log.Printf("Cleaning Temp Directory: %s\n", mockFsRoot)
		err := os.RemoveAll(mockFsRoot)
		if err != nil {
			log.Printf("Error cleaning Temp (%s): %s\n", mockFsRoot, err)
		}
	}

	tree, err := NewTree(objects, mockFsRoot)
	if err != nil {
		ReloadMainBucketsWithError(err)
	}

	listing := GetNodeDirectory(tree)
	if err != nil {
		ReloadMainBucketsWithError(err)
	}

	SetBackHandler(ReloadMainBuckets, deferFunc)
	RenderBucketExplorerListing(bucket, listing, selection, deferFunc)

}
