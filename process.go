package main

import (
	"strings"
)

func extract(content string) ([]wordNode, *parseError) {
	var group []wordNode

	currentContent := content
	defineIndex := strings.Index(currentContent, "define ")
	if defineIndex != -1 {
		firstBracketIndex := strings.Index(currentContent, "{")
		lastBracketIndex := strings.Index(currentContent, "}")

		var name string
		var defType string
		var node wordNode

		defineField := strings.TrimSpace(currentContent[(defineIndex + 7):(firstBracketIndex)])
		defineGroup := strings.Split(defineField, " ")
		if len(defineGroup) >= 2 {
			defType = defineGroup[0]
			name = defineGroup[1]
		} else {
			pe := parseError{
				Period: "extract",
				Description: "Wrong syntax near \"define\"",
			}

			return group, &pe
		}

		// Validate
		nameFlag := validateName(name)
		typeFlag := validateType(defType)
		if !nameFlag {
			pe := parseError{
				Period: "extract",
				Description: "Invalid name \"" + name + "\"",
			}

			return group, &pe
		} else if !typeFlag {
			pe := parseError{
				Period: "extract",
				Description: "Invalid type \"" + defType + "\"",
			}

			return group, &pe
		}

		// Initialize node
		attrGroup := [0] wordAttr {}
		node = wordNode{
			Name: name,
			TypeName: defType,
			AttrGroup: attrGroup[:],
		}

		// Extract Attribute
		defineContent := currentContent[(firstBracketIndex + 1):lastBracketIndex]
		contentGroup := strings.Split(defineContent, ",")
		for _, x := range contentGroup {
			var attr wordAttr

			splitGroup := strings.Split(x, ":")
			if len(splitGroup) == 2 {
				left := strings.TrimSpace(splitGroup[0])
				left = strings.Replace(left, "/r", "", -1)
				left = strings.Replace(left, "/n", "", -1)
				left = strings.Replace(left, "/r/n", "", -1)
				right := strings.TrimSpace(splitGroup[1])

				var searchGroup []string
				if defType == "exp" {
					searchGroup = expRequired[:]
				} else if defType == "time" {
					searchGroup = timeRequired[:]
				} else if defType == "place" {
					searchGroup = placeRequired[:]
				} else if defType == "role" {
					searchGroup = roleRequired[:]
				}

				for _, x := range searchGroup {
					if left == x || ("@" + left) == x {
						attr = wordAttr{
							Value: right,
							AttrName: left,
						}

						node.AttrGroup = append(node.AttrGroup, attr)
					}
				}

			} else {
				pe := parseError{
					Period: "extract",
					Description: "Wrong syntax in attribute definition",
				}

				return group, &pe
			}

			group = append(group, node)
		}
	}

	return group, nil
}

func analyse(group []wordNode) ([]wordNode, expNode) {
	var words []wordNode
	var exp expNode

	return words, exp
}

func explain() {

}
