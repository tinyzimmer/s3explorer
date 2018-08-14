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
	"os"

	"github.com/gizak/termui"
)

func RenderBucketListing(buckets []BucketWithDisplay) {

	// Set the selection pointer

	var selection int
	selection = 0

	// Create a UI ready list and render

	list := CreateBucketList(buckets, selection)
	termui.Render(list, RenderHelp())

	// up goes up

	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		if selection == 0 {
			return
		} else {
			selection -= 1
			list := CreateBucketList(buckets, selection)
			termui.Render(list, RenderHelp())
		}
	})

	// down goes sideways

	termui.Handle("/sys/kbd/<down>", func(termui.Event) {
		if selection == len(buckets)-1 {
			return
		} else {
			selection += 1
			list := CreateBucketList(buckets, selection)
			termui.Render(list, RenderHelp())
		}
	})

	// enter loads the bucket

	termui.Handle("/sys/kbd/<enter>", func(termui.Event) {
		termui.ResetHandlers()
		p := RenderMessage("Loading Bucket", buckets[selection].displayString)
		termui.Render(p)
		RenderBucketExplorer(buckets[selection])
	})

}

func ReloadMainBucketsWithError(err error) {

	// Render an error to the screen and reload the bucket listing

	RenderError(err.Error())
	ReloadMainBuckets()
}

func ReloadMainBuckets() {

	// Retrieve a fresh bucket listing and reload main screen

	termui.ResetHandlers()
	SetDefaultHandlers(func() { return })
	termui.Clear()
	buckets, err := s3Session.GetBucketWithDisplayStrings()
	if err != nil {
		os.Exit(EXIT_FAILED_BUCKET_LISTING)
	}
	RenderBucketListing(buckets)
}
