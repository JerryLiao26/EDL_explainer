package main

var superWord = [...]string{"define"}
var typeWord = [...]string{"exp", "time", "place", "role"}
var keyWord = [...]string{"time", "place", "role", "process", "input", "output", "addr", "link"}
var preWord = [...]string{"any", "start", "end", "if", "elif", "else"}

var expRequired = [...]string{"time", "place", "role", "process", "@input", "@output"}
var timeRequired = [...]string{"@period", "@exact", "@last"}
var placeRequired = [...]string{"@addr", "@link"}
var roleRequired = [...]string{"title"}

type wordNode struct {
	Name      string     `json:"name"`
	TypeName  string     `json:"type_name"`
	AttrGroup []wordAttr `json:"attr_group"`
}

type wordAttr struct {
	Value    string `json:"value"`
	AttrName string `json:"attr_name"`
}

type expNode struct {
	Time    wordNode `json:"time"`
	Role    wordNode `json:"role"`
	Place   wordNode `json:"place"`
	Process string   `json:"process"`
	Name    string   `json:"name"`
}

type serverRespond struct {
	Code   int    `json:"code"`
	Text   string `json:"text"`
	Method string `json:"method"`
	Status bool   `json:"status"`
}

type clientMessage struct {
	Target  string `json:"target"`
	Content string `json:"content"`
}
