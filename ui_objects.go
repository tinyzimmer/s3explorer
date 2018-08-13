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
	"strings"
	"time"

	"github.com/gizak/termui"
)

func RenderHelp() (p *termui.Par) {
	arrows := "\u2195\ufe0f"
	returnArrow := "\u21b2"
	helpText := fmt.Sprintf("%v navigate - %v open - <q> quit - <b> back", arrows, returnArrow)
	p = termui.NewPar(helpText)
	p.Height = 3
	p.Width = len(helpText) + 3
	p.TextFgColor = termui.ColorWhite
	p.BorderLabel = "Help"
	p.BorderFg = termui.ColorCyan
	p.Y = termui.TermHeight() - 5
	return
}

func RenderMessage(label string, message string) (p *termui.Par) {
	p = termui.NewPar(message)
	p.Height = 3
	p.Width = termui.TermWidth() - RIGHT_BUFFER
	p.TextFgColor = termui.ColorWhite
	p.BorderLabel = label
	p.BorderFg = termui.ColorCyan
	return
}

func RenderError(errorMessage string) {
	p := termui.NewPar(errorMessage)
	p.Height = 3
	p.Width = termui.TermWidth() - RIGHT_BUFFER
	p.TextFgColor = termui.ColorWhite
	p.BorderLabel = "Error"
	p.BorderFg = termui.ColorCyan
	termui.Render(p)
	time.Sleep(time.Duration(time.Second * 2))
}

func CreateDownloadPrompt(node *Node) (p *termui.Par) {
	p = termui.NewPar(*node.S3Object.Key)
	p.Height = 5
	p.Width = termui.TermWidth() - RIGHT_BUFFER
	p.TextFgColor = termui.ColorWhite
	p.Border = false
	p.Y = termui.TermHeight() - 10
	return
}

func GetBucketListing(buckets []BucketWithDisplay, selection int) (listing []string) {
	var index int
	index = 0
	for _, bucket := range buckets {
		if index == selection {
			listing = append(listing, fmt.Sprintf("[[%v] %s](bg-blue)", index, bucket.displayString))
		} else {
			listing = append(listing, fmt.Sprintf("[%v] %s", index, bucket.displayString))
		}
		index += 1
	}
	return
}

func CreateBucketList(buckets []BucketWithDisplay, selection int) *termui.List {

	var displayStrings []string

	for _, bucket := range buckets {
		displayStrings = append(displayStrings, bucket.displayString)
	}

	ls := termui.NewList()
	ls.Items = GetBucketListing(buckets, selection)
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = "S3 Buckets"
	ls.Height = (len(buckets) + 2)
	ls.Width = termui.TermWidth() - RIGHT_BUFFER
	ls.Y = 0
	return ls
}

func TruncateFilename(filename string) (truncated string, space int) {
	if len(filename) >= termui.TermWidth()/4 {
		truncated = fmt.Sprintf("%s...", filename[:(termui.TermWidth()/4)-3])
	} else {
		truncated = filename
	}
	space = (termui.TermWidth() / 2) - len(truncated)
	return
}

func GetDirectoryDisplayListing(objects []string, selection int) (listing []string) {
	var index int
	index = 0
	for _, obj := range objects {
		if index == selection {
			listing = append(listing, fmt.Sprintf("[[%v] %s](bg-blue)", index, obj))
		} else {
			listing = append(listing, fmt.Sprintf("[%v] %s", index, obj))
		}
		index += 1
	}
	return
}

func CreateDirectoryList(title string, nodes []*Node, selection int) *termui.List {

	var displayStrings []string

	for _, node := range nodes {
		var display string
		if !node.Info.IsDir {
			file, space := TruncateFilename(node.DisplayString)
			display = fmt.Sprintf("%s%s%v", file, strings.Repeat(" ", space), ByteFormat(float64(*node.S3Object.Size), 1))
		} else {
			display = node.DisplayString
		}
		displayStrings = append(displayStrings, display)
	}

	ls := termui.NewList()
	ls.Items = GetDirectoryDisplayListing(displayStrings, selection)
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = title
	ls.Height = (len(nodes) + 2)
	ls.Width = termui.TermWidth() - RIGHT_BUFFER
	ls.Y = 0
	return ls
}
