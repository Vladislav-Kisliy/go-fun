package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	// "strconv"
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
	// if !(len(os.Args) == 2 || len(os.Args) == 3) {
	// 	panic("usage go run main.go . [-f]")
	// }
	// path := os.Args[1]
	// path := "/home/vlad/work/projects/school/go/course/hw1_tree"
	path := "testdata"
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
	fmt.Println("out =", fileRoot)
	fmt.Println("out nodes =", fileRoot.nodes)

	// var files []string
	// for _, file := range fileInfos {
	// 	// files = append(files, file.Name())
	// 	fmt.Fprintln(out, file)
	// }

	return nil
}

func readFolder(root string, printFiles bool) (FileNode, error) {
	fileRoot := FileNode{Filepath: root, fileInfo: FileInfo{isDir: true}}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, separator) {
			var shouldAdd = false
			filepath := transformPath(root, path)
			lastIndex := len(filepath) - 1
			if lastIndex >= 0 {
				fmt.Println("name will be " + filepath[lastIndex])
				fileNode := FileNode{
					Filepath: filepath[lastIndex],
				}
				fileInfo := FileInfo{}
				if info.IsDir() {
					fileInfo.isDir = true
					shouldAdd = true
					// fileInfos = append(fileInfos, transformPath(root, path))
				} else if printFiles {
					fileInfo.Size = info.Size()
					shouldAdd = true
					// fileInfo = append(fileInfo, transformPath(root, path)+" ("+string(info.Size())+")")
					// if info.Size() == 0 {
					// 	fileInfos = append(fileInfos, transformPath(root, path)+" (empty)")
					// } else {
					// 	fileInfos = append(fileInfos, transformPath(root, path)+" ("+strconv.FormatInt(info.Size(), 10)+"b)")
					// }
				}
				if shouldAdd {
					fileNode.fileInfo = fileInfo
					insertNode(&fileRoot, fileNode, filepath)
					// fileRoot.nodes = append(fileRoot.nodes, fileNode)
					shouldAdd = false
					fmt.Println("fileNode =", fileNode)
				}
			}

		}

		return nil
	})
	return fileRoot, err
}

func transformPath(root string, path string) []string {
	var result []string
	fmt.Println("path =", path)
	if indexRoot := strings.Index(path, root); indexRoot > -1 {
		clearLine := path[indexRoot:]
		fmt.Println("clear data ", clearLine)
		// clearLine := strings.Replace(path, root+separator, "", 1)
		result = strings.Split(clearLine, separator)
	}

	return result
}

func insertNode(rootNode *FileNode, targetNode FileNode, path []string) bool {
	pathLen := len(path)
	fmt.Println("path =", path)
	if rootNode.fileInfo.isDir && pathLen > 0 && rootNode.Filepath == path[0] {

		switch pathLen {
		case 1:
			fmt.Println("Target 1", path)
			rootNode.nodes = append(rootNode.nodes, targetNode)

		case 2:
			fmt.Println("Target 2", path)

			fmt.Println("Target 2 add", targetNode)
			rootNode.nodes = append(rootNode.nodes, targetNode)
			// fileRoot.nodes = append(fileRoot.nodes, fileNode)
			fmt.Println("Target 2 nodes after", rootNode.nodes)
			return true

		default:
			fmt.Println("Target >2", path)
			for i := 0; i < len(rootNode.nodes); i++ {
				if insertNode(&rootNode.nodes[i], targetNode, path[1:]) {
					return true
				}
			}
		}
	}
	return false
}

func genLine(parentNode FileNode, position int, firstSymbol string) string {
	nodesAmount := len(parentNode.nodes)
	var indent string
	if firstSymbol == lastElement {
		indent = strings.Repeat("\t", position+1)
	} else if firstSymbol == middleElement {
		indent = verticalLine + strings.Repeat("\t", position)

	}

	var result string
	if nodesAmount > 0 {
		for i := 0; i < nodesAmount; i++ {
			currentElement := middleElement
			if i == (nodesAmount - 1) {
				currentElement = lastElement
			}
			result = result + "\n" + indent + genLine(parentNode.nodes[i], position+1, currentElement)
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
	result := "\n"
	nodesAmount := len(fileNode.nodes)
	if nodesAmount > 0 {
		for i := 0; i < nodesAmount; i++ {
			currentElement := middleElement
			if i == (nodesAmount - 1) {
				currentElement = lastElement
			}
			result = result + genLine(fileNode.nodes[i], 0, currentElement) + "\n"
		}
	}

	return result
}
