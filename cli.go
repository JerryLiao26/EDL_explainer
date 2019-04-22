package main

import (
	"fmt"
)

func explainCliHandler() {
	fmt.Println("Explain")
}

func analyseCliHandler(path string) {
	fullPath := getAbsPath(path)

	stream := readLocalFile(fullPath)
	group, pe := extract(string(stream))

	if pe != nil {
		fmt.Println(stringifyParseError(pe))
	} else {
		nodes, pe := analyse(group)
		if pe != nil {
			fmt.Println(stringifyParseError(pe))
		} else {
			writeLocalFile(fullPath + ".ana", []byte(stringifyExpNodes(nodes)))
		}
	}
}

func extractCliHandler(path string) {
	fullPath := getAbsPath(path)

	stream := readLocalFile(fullPath)
	group, pe := extract(string(stream))

	if pe != nil {
		fmt.Println(stringifyParseError(pe))
	} else {
		writeLocalFile(fullPath + ".ext", []byte(stringifyWordNodes(group)))
	}
}
