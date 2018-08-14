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
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
)

type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
}

func fileInfoFromInterface(v os.FileInfo) (info *FileInfo, display string) {
	info = &FileInfo{v.Name(), v.Size(), v.Mode(), v.ModTime(), v.IsDir()}
	if info.IsDir {
		display = info.Name + localDelimiter
	} else {
		display = info.Name
	}
	return
}

type Node struct {
	FullPath      string
	DisplayString string
	Info          *FileInfo
	Children      []*Node
	Parent        *Node
	S3Object      *s3.Object
}

func CreateMockFs(objects []*s3.Object) (tempDir string, err error) {
	log.Println("Creating Mock filesystem for bucket indexing")
	tempDir, err = ioutil.TempDir("", "")
	if err != nil {
		return
	}
	for _, obj := range objects {
		dir, file := path.Split(*obj.Key)
		sanitizedDir := strings.Replace(dir, "/", localDelimiter, -1)
		if len(dir) > 0 {
			err := os.MkdirAll(path.Join(tempDir, sanitizedDir), DEFAULT_DIRECTORY_MODE)
			if err != nil {
				log.Printf("Error: %s\n", err.Error())
			}
		}
		if len(file) > 0 {
			sysFile, err := os.Create(path.Join(tempDir, sanitizedDir, file))
			if err != nil {
				log.Printf("Error: %s\n", err.Error())
			}
			sysFile.Close()
		}
	}
	log.Println("Created mock filesystem")
	return
}

func MatchS3Object(objects []*s3.Object, root string, path string) *s3.Object {
	sanitized := strings.Replace(strings.Replace(path, root+localDelimiter, "", 1), localDelimiter, "/", -1)
	for _, obj := range objects {
		if *obj.Key == sanitized {
			return obj
		}
	}
	return nil
}

// Create directory hierarchy.
func NewTree(objects []*s3.Object, root string) (result *Node, err error) {
	log.Printf("Creating node tree for bucket filesystem at: %s\n", root)
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return
	}
	parents := make(map[string]*Node)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		cInfo, display := fileInfoFromInterface(info)
		var objectMatch *s3.Object
		if len(objects) > 0 {
			objectMatch = MatchS3Object(objects, root, path)
		} else {
			objectMatch = nil
		}
		parents[path] = &Node{
			FullPath:      path,
			DisplayString: display,
			Info:          cInfo,
			Children:      make([]*Node, 0),
			S3Object:      objectMatch,
		}
		return nil
	}
	if err = filepath.Walk(absRoot, walkFunc); err != nil {
		return
	}
	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists { // If a parent does not exist, this is the root.
			result = node
		} else {
			node.Parent = parent
			parent.Children = append(parent.Children, node)
		}
	}
	log.Println("Finished indexing bucket fs")
	return
}

func GetLocalDelimiter() string {
	log.Println("Checking local delimiter")
	if runtime.GOOS == "windows" {
		return "\\"
	} else {
		return "/"
	}
}

func GetSubdirs(node *Node) (nodes []*Node) {
	log.Printf("Getting subdirs for node: %+v\n", node)
	for _, child := range node.Children {
		if child.Info.IsDir {
			nodes = append(nodes, child)
		}
	}
	return
}

func GetFiles(node *Node) (nodes []*Node) {
	log.Printf("Getting files for node: %+v\n", node)
	for _, child := range node.Children {
		if !child.Info.IsDir {
			nodes = append(nodes, child)
		}
	}
	return
}

func GetNodeDirectory(node *Node) (nodes []*Node) {
	log.Printf("Creating node directory tree for node focus: %+v\n", node)
	if node.Parent != nil {
		nodes = append(nodes, &Node{
			FullPath:      node.Parent.FullPath,
			DisplayString: "..",
			Info:          node.Parent.Info,
			Children:      node.Parent.Children,
			Parent:        node.Parent.Parent,
			S3Object:      node.Parent.S3Object,
		})
	}
	nodes = append(nodes, GetSubdirs(node)...)
	nodes = append(nodes, GetFiles(node)...)
	return
}
