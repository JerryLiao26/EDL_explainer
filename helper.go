package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

func readLocalFile(path string) []byte {
	b, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println("error:", err)
	}

	return b
}

func writeLocalFile(path string, content []byte) {
	err := ioutil.WriteFile(path, content, 0644)

	if err != nil {
		fmt.Println("error:", err)
	}
}

func getAbsPath(path string) string {
	fullPath, err := filepath.Abs(path)

	if err != nil {
		fmt.Println("error:", err)
	}

	return fullPath
}

func methodNotAllow(w http.ResponseWriter, r *http.Request) {
	sr := serverRespond{
		Code:   405,
		Status: false,
		Method: r.Method,
		Text:   r.Method + " method not allowed",
	}

	stream, err := json.Marshal(sr)
	if err != nil {
		fmt.Println("error:", err)
	}

	_, _ = fmt.Fprintf(w, string(stream))
}

func validateName(name string) bool {
	var preservedWords []string
	var prefixedName = "@" + name

	// Illegal Characters
	for i := 0; i <= 126; i++ {
		if (i >= 65 && i <= 90) || i == 95 || (i >= 97 && i <= 122) {
			continue
		}

		rIndex := strings.IndexRune(name, rune(i))

		if i >= 48 && i <= 57 && rIndex != -1 { // Cannot start with
			return false
		} else if rIndex != -1 { // Cannot contain
			return false
		}
	}

	// Preserved words
	preservedWords = append(superWord[:], typeWord[:]...)
	preservedWords = append(preservedWords, keyWord[:]...)
	preservedWords = append(preservedWords, preWord[:]...)

	joinedString := strings.Join(preservedWords, ",")

	if strings.Contains(joinedString, name) || strings.Contains(joinedString, prefixedName) {
		return false
	}

	// Passed
	return true
}

func validateType(typeName string) bool {
	typeString := strings.Join(typeWord[:], ",")

	if strings.Contains(typeString, typeName) {
		return true
	}

	return false
}

func stringifyParseError(pe *parseError) string {
	return "Error during " + pe.Period + ": " + pe.Description
}

func stringifyWordNodes(group []wordNode) string {
	var jsonString = "["

	for _, x := range group {
		stream, err := json.Marshal(x)

		if err != nil {
			fmt.Println("error:", err)
		} else {
			jsonString += string(stream) + ","
		}
	}

	println(jsonString)
	jsonString = jsonString[:(len(jsonString) - 1)] + "]"

	return jsonString
}
