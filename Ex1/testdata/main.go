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
	lastElement   = "└"
	middleElement = "├"
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
	fmt.Println("out ", fileRoot)

	// var files []string
	// for _, file := range fileInfos {
	// 	// files = append(files, file.Name())
	// 	fmt.Fprintln(out, file)
	// }

	return nil
}

func readFolder(root string, printFiles bool) (FileNode, error) {

	fileRoot := FileNode{Filepath: root}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, separator) {
			var shouldAdd = false
			fileNode := FileNode{
				Filepath: transformPath(root, path),
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
				// fileInfos = append(fileInfos, fileInfo)
				fileNode.fileInfo = fileInfo
				fileRoot.nodes = append(fileRoot.nodes, fileNode)
				shouldAdd = false
				fmt.Println("fileNode =", fileNode)
			}

		}

		return nil
	})
	return fileRoot, err
}
func transformPath(root string, path string) []string {

	fmt.Println("path =", path)
	clearLine := strings.Replace(path, root+separator, "", 1)
	result := strings.Split(clearLine, separator)

	// if lenParts > 0 {
	// 	result = strings.Replace(path, root+separator, "", 1)
	// 	// strings.LastIndex(path, separator)
	// 	// nameIndex := strings.LastIndex(path, separator) + 1
	// 	// fmt.Println("index =", nameIndex)
	// 	// switch lenParts {
	// 	// case 2:
	// 	// 	result = strings.Replace(path, root+separator, "", 1)
	// 	// case 3:
	// 	// 	result = path[nameIndex:]
	// 	// case 4:
	// 	// 	result = path[nameIndex:]
	// 	// }
	// }

	return result
}

func genLine(fileNode FileNode, position int, firstSymbol string) string {
	if len(fileNode.nodes) > 1 {
		for i := 0; i < len(fileNode.nodes)-2; i++ {
			if fileNode.nodes[i].fileInfo.isDir && len(fileNode.nodes[i].nodes) > 0 {
				genLine(fileNode.nodes[i], position+1, middleElement)
			}
		}
	}
	if firstSymbol == lastElement {
		return "last" + firstSymbol + fileNode.Filepath
	} else if firstSymbol == middleElement {
		return firstSymbol + fileNode.Filepath
	}
	return ""
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
	var result string
	for _, innerNode := range fileNode.nodes {
		result = result + genLine(innerNode, 1, middleElement) + "\n"
	}
	return result
	// return "├───" + fileNode.Filepath + "nodes []"
}
