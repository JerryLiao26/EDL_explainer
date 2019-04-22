package main

import (
	"strings"
)

func extract(content string) ([]wordNode, *parseError) {
	var group []wordNode

	allNames := ""

	// Replace break-lines
	currentContent := strings.Replace(content, "/r", "", -1)
	currentContent = strings.Replace(currentContent, "/n", "", -1)
	currentContent = strings.Replace(currentContent, "/r/n", "", -1)
	defineIndex := strings.Index(currentContent, "define ")
	for defineIndex != -1 {
		firstBracketIndex := strings.Index(currentContent, "{")
		lastBracketIndex := strings.Index(currentContent, "}")
		if firstBracketIndex == -1 || lastBracketIndex == -1 {
			pe := parseError{
				Period:      "extract",
				Description: "Missing define bracket",
			}

			return group, &pe
		}

		var name string
		var defType string
		var node wordNode

		defineField := strings.TrimSpace(currentContent[(defineIndex + 7):(firstBracketIndex)])
		defineGroup := strings.Split(defineField, " ")
		if len(defineGroup) >= 2 {
			defType = defineGroup[0]
			name = defineGroup[1]

			// Check duplicate name
			if fullyContains(allNames, name) {
				pe := parseError{
					Period:      "extract",
					Description: "Duplicate name \"" + name + "\"",
				}

				return group, &pe
			}

			// For fullyContains helper function
			if allNames == "" {
				allNames += name
			} else {
				allNames += "," + name
			}
		} else {
			pe := parseError{
				Period:      "extract",
				Description: "Wrong syntax near \"define\"",
			}

			return group, &pe
		}

		// Validate
		nameFlag := validateName(name)
		typeFlag := validateType(defType)
		if !nameFlag {
			pe := parseError{
				Period:      "extract",
				Description: "Invalid name \"" + name + "\"",
			}

			return group, &pe
		} else if !typeFlag {
			pe := parseError{
				Period:      "extract",
				Description: "Invalid type \"" + defType + "\"",
			}

			return group, &pe
		}

		// Initialize node
		attrGroup := [0]wordAttr{}
		node = wordNode{
			Name:      name,
			TypeName:  defType,
			AttrGroup: attrGroup[:],
		}

		// Extract Attribute
		defineContent := currentContent[(firstBracketIndex + 1):lastBracketIndex]
		contentGroup := strings.Split(defineContent, ",")
		for _, x := range contentGroup {
			var attr wordAttr

			splitGroup := strings.SplitN(x, ":", 2)
			if len(splitGroup) == 2 {
				left := strings.TrimSpace(splitGroup[0])
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
				} else {
					pe := parseError{
						Period:      "extract",
						Description: "Unknown type \"" + defType + "\"",
					}

					return group, &pe
				}

				for _, x := range searchGroup {
					if left == x || ("@"+left) == x {
						attr = wordAttr{
							Value:    right,
							AttrName: left,
						}

						node.AttrGroup = append(node.AttrGroup, attr)

						// Found, pass
						break
					}
				}

			} else {
				pe := parseError{
					Period:      "extract",
					Description: "Wrong syntax in attribute definition",
				}

				return group, &pe
			}
		}
		group = append(group, node)

		// Reset loop values
		currentContent = currentContent[(lastBracketIndex + 1):]
		defineIndex = strings.Index(currentContent, "define ")
	}
	return group, nil
}

func analyse(group []wordNode) ([]expNode, *parseError) {
	var nodes []expNode
	var searchGroup []string

	var expNames []string

	var timeNames []string
	timeMap := make(map[string]int)

	var placeNames []string
	placeMap := make(map[string]int)

	var roleNames []string
	roleMap := make(map[string]int)

	// Check required attr
	for i, x := range group {
		if x.TypeName == "exp" { // Each word
			searchGroup = expRequired[:]
			expNames = append(expNames, x.Name)
		} else if x.TypeName == "time" {
			searchGroup = timeRequired[:]
			timeNames = append(timeNames, x.Name)
			timeMap[x.Name] = i
		} else if x.TypeName == "place" {
			searchGroup = placeRequired[:]
			placeNames = append(placeNames, x.Name)
			placeMap[x.Name] = i
		} else if x.TypeName == "role" {
			searchGroup = roleRequired[:]
			roleNames = append(roleNames, x.Name)
			roleMap[x.Name] = i
		}

		ownsAttrString := ""
		for _, xx := range x.AttrGroup { // Each attr of word
			ownsAttrString += xx.AttrName + ","
		}

		if ownsAttrString != "" {
			for _, xx := range searchGroup {
				if !strings.HasPrefix(xx, "@") { // Optional
					if !fullyContains(ownsAttrString, xx) {
						pe := parseError{
							Period:      "analyse",
							Description: "Missing attribute \"" + xx + "\" for " + x.TypeName + " \"" + x.Name + "\"",
						}

						return nodes, &pe
					}
				}
			}
		} else {
			pe := parseError{
				Period:      "analyse",
				Description: "Must have at least one attribute for " + x.TypeName + " \"" + x.Name + "\"",
			}

			return nodes, &pe
		}
	}

	// Check attribute content
	expNamesString := strings.Join(expNames, ",")
	timeNamesString := strings.Join(timeNames, ",")
	roleNamesString := strings.Join(roleNames, ",")
	placeNamesString := strings.Join(placeNames, ",")

	for _, x := range group {
		var exp expNode
		var currentStep *stepNode

		// Initialize exp node
		if x.TypeName == "exp" {
			exp = expNode{
				Name: x.Name,
				Process: nil,
			}
		}

		for _, xx := range x.AttrGroup {
			if xx.AttrName == "time" || xx.AttrName == "place" || xx.AttrName == "role" {
				if xx.Value == "any" {
					continue // Passed
				} else if xx.AttrName == "time" && fullyContains(timeNamesString, xx.Value) {
					exp.Time = group[timeMap[xx.Value]]
					continue
				} else if xx.AttrName == "place" && fullyContains(placeNamesString, xx.Value) {
					exp.Place = group[placeMap[xx.Value]]
					continue
				} else if xx.AttrName == "role" && fullyContains(roleNamesString, xx.Value) {
					exp.Role = group[roleMap[xx.Value]]
					continue
				} else if wrapByDoubleQuotes(xx.Value) {
					continue
				}
			} else if xx.AttrName == "process" {
				processGroup := strings.Split(xx.Value, "=>")
				if strings.TrimSpace(processGroup[0]) != "start" || strings.TrimSpace(processGroup[len(processGroup)-1]) != "end" {
					pe := parseError{
						Period:      "analyse",
						Description: "Process of exp \"" + x.Name + "\" should start with \"start\" and end with \"end\"",
					}

					return nodes, &pe
				}

				// Check step content
				for i := 1; i < len(processGroup)-1; i++ {
					var step stepNode

					stepContent := strings.TrimSpace(processGroup[i])
					if wrapByDoubleQuotes(stepContent) {
						branches := [0]*stepNode{}
						step = stepNode{
							Name: unwrap(stepContent),
							TypeName: "",
							Condition: "",
							DirectChild: nil,
							Branches: branches[:],
						}
					} else if fullyContains(expNamesString, stepContent) {
						branches := [0]*stepNode{}
						step = stepNode{
							Name: stepContent,
							TypeName: "exp",
							Condition: "",
							DirectChild: nil,
							Branches: branches[:],
						}
					} else if wrapBySquareBrackets(stepContent) {
						currentStepContent := stepContent
						leftBracketIndex := strings.Index(currentStepContent, "(")
						rightBracketIndex := strings.Index(currentStepContent, ")")
						for leftBracketIndex != -1 {
							var subStep stepNode

							// Bracket not matched
							if rightBracketIndex == -1 {
								pe := parseError{
									Period:      "analyse",
									Description: "A \"(\" in process of exp \"" + x.Name + "\" is not matched",
								}

								return nodes, &pe
							}

							conditionContent := currentStepContent[(leftBracketIndex + 1):rightBracketIndex]
							ifIndex := strings.Index(conditionContent, "if")
							doIndex := strings.Index(conditionContent, "do")
							if ifIndex == -1 || doIndex == -1 {
								pe := parseError{
									Period:      "analyse",
									Description: "Condition in process of exp \"" + x.Name + "\" doesn't match if...do... clause",
								}

								return nodes, &pe
							}

							// Condition&Job
							condition := strings.TrimSpace(conditionContent[(ifIndex + 2):(doIndex - 1)])
							job := strings.TrimSpace(conditionContent[(doIndex + 2):])

							if !wrapByDoubleQuotes(condition) {
								pe := parseError{
									Period:      "analyse",
									Description: "Condition if statement in process of exp \"" + x.Name + "\" should be wrap by double quotes",
								}

								return nodes, &pe
							} else {
								if wrapByDoubleQuotes(job) {
									branches := [0]*stepNode{}
									subStep = stepNode{
										Name: unwrap(job),
										TypeName: "",
										Condition: unwrap(condition),
										DirectChild: nil,
										Branches: branches[:],
									}
								} else if fullyContains(expNamesString, job) && job != x.Name {
									branches := [0]*stepNode{}
									subStep = stepNode{
										Name: job,
										TypeName: "exp",
										Condition: unwrap(condition),
										DirectChild: nil,
										Branches: branches[:],
									}
								} else {
									pe := parseError{
										Period:      "analyse",
										Description: "Condition do statement in process of exp \"" + x.Name + "\" should be wrap by double quotes or use another exp name",
									}

									return nodes, &pe
								}

								// Append to process steps
								if currentStep.DirectChild == nil {
									currentStep.DirectChild = &subStep
									step = subStep
								} else {
									currentStep.Branches = append(currentStep.Branches, &subStep)
								}
							}

							// Reset loop conditions
							currentStepContent = currentStepContent[(rightBracketIndex + 1):]
							leftBracketIndex = strings.Index(currentStepContent, "(")
							rightBracketIndex = strings.Index(currentStepContent, ")")
						}
					} else {
						pe := parseError{
							Period:      "analyse",
							Description: "Process of exp \"" + x.Name + "\" contains bad definition",
						}

						return nodes, &pe
					}

					// Append to exp
					if exp.Process == nil {
						exp.Process = &step
						currentStep = &step
					} else {
						currentStep.DirectChild = &step
						currentStep = &step
					}
				}

				// Passed
				continue
			} else if wrapByDoubleQuotes(xx.Value) {
				continue
			}

			// Not passed
			pe := parseError{
				Period:      "analyse",
				Description: "Wrong value for content of \"" + xx.AttrName + "\" in exp \"" + x.Name + "\"",
			}

			return nodes, &pe
		}

		// Append to nodes
		if x.TypeName == "exp" {
			nodes = append(nodes, exp)
		}
	}

	return nodes, nil
}

func explain(nodes []expNode, target string) string {
	return ""
}
