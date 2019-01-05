package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	separator     = string(os.PathSeparator)
	lastElement   = "└───"
	middleElement = "├───"
	verticalLine  = "│\t"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	// path := "testdata"
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	// printFiles := true
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, root string, printFiles bool) error {
	fileRoot, err := readFolder(root, printFiles)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// fmt.Println("out =", fileRoot)
	// fmt.Println("out nodes =", fileRoot.nodes)

	fmt.Fprint(out, fileRoot)

	return err
}

func readFolder(root string, printFiles bool) (FileNode, error) {
	fileRoot := FileNode{Filepath: root, fileInfo: FileInfo{isDir: true}}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, separator) {
			var shouldAdd = false
			filepath := transformPath(root, path)
			lastIndex := len(filepath) - 1
			if lastIndex >= 0 {
				fileNode := FileNode{
					Filepath: filepath[lastIndex],
				}
				fileInfo := FileInfo{}
				if info.IsDir() {
					fileInfo.isDir = true
					shouldAdd = true
				} else if printFiles {
					fileInfo.Size = info.Size()
					shouldAdd = true
					if info.Size() == 0 {
						fileNode.Filepath = fileNode.Filepath + " (empty)"
					} else {
						fileNode.Filepath = fileNode.Filepath + " (" + strconv.FormatInt(info.Size(), 10) + "b)"
					}
				}
				if shouldAdd {
					fileNode.fileInfo = fileInfo
					insertNode(&fileRoot, fileNode, filepath)
					shouldAdd = false
				}
			}

		}

		return nil
	})
	return fileRoot, err
}

func transformPath(root string, path string) []string {
	var result []string
	if indexRoot := strings.Index(path, root); indexRoot > -1 {
		clearLine := path[indexRoot:]
		result = strings.Split(clearLine, separator)
	}

	return result
}

func insertNode(rootNode *FileNode, targetNode FileNode, path []string) bool {
	pathLen := len(path)
	// fmt.Println("path =", path)
	if rootNode.fileInfo.isDir && pathLen > 0 && rootNode.Filepath == path[0] {
		switch pathLen {
		case 1:
			rootNode.nodes = append(rootNode.nodes, targetNode)

		case 2:
			rootNode.nodes = append(rootNode.nodes, targetNode)
			return true

		default:
			for i := 0; i < len(rootNode.nodes); i++ {
				if insertNode(&rootNode.nodes[i], targetNode, path[1:]) {
					return true
				}
			}
		}
	}
	return false
}

func genLine(parentNode FileNode, position int, firstSymbol string, isLastElement bool) string {
	nodesAmount := len(parentNode.nodes)
	var indent string
	if firstSymbol == lastElement {
		if isLastElement {
			indent = strings.Repeat("\t", position+1)
		} else {
			indent = strings.Repeat(verticalLine, position) + "\t"
		}

	} else if firstSymbol == middleElement {
		indent = strings.Repeat(verticalLine, position+1)
	}

	var result string
	if nodesAmount > 0 {
		for i := 0; i < nodesAmount; i++ {
			currentElement := middleElement
			if i == (nodesAmount - 1) {
				currentElement = lastElement
			}
			result = result + "\n" + indent + genLine(parentNode.nodes[i], position+1, currentElement, isLastElement)
		}
	}

	result = firstSymbol + parentNode.Filepath + result

	return result
}

type FileInfo struct {
	Size  int64
	isDir bool
}

type FileNode struct {
	Filepath string
	fileInfo FileInfo
	nodes    []FileNode
}

func (fileNode FileNode) String() string {
	result := ""
	nodesAmount := len(fileNode.nodes)
	if nodesAmount > 0 {
		for i := 0; i < nodesAmount-1; i++ {
			result = result + genLine(fileNode.nodes[i], 0, middleElement, false) + "\n"
		}
		result = result + genLine(fileNode.nodes[nodesAmount-1], 0, lastElement, true) + "\n"
	}

	return result
}
