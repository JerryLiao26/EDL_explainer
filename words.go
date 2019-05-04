package main

var superWord = [...]string{"define"}
var typeWord = [...]string{"exp", "time", "place", "role"}
var keyWord = [...]string{"time", "place", "role", "process", "material", "input", "output", "period", "exact", "last", "addr", "link", "desc"}
var preWord = [...]string{"any", "start", "end", "if", "do"}

var expRequired = [...]string{"time", "place", "role", "process", "@material", "@input", "@output"}
var timeRequired = [...]string{"@period", "@exact", "@last"}
var placeRequired = [...]string{"@desc", "@addr", "@link"}
var roleRequired = [...]string{"desc"}

type parseError struct {
	Period      string `json:"period"`
	Description string `json:"desc"`
}

type wordNode struct {
	Name      string     `json:"name"`
	TypeName  string     `json:"type_name"`
	AttrGroup []wordAttr `json:"attr_group"`
}

type wordAttr struct {
	Value    string `json:"value"`
	AttrName string `json:"attr_name"`
}

type stepNode struct {
	Name        string      `json:"name"`
	TypeName    string      `json:"type"`
	Condition   string      `json:"condition"`
	DirectChild *stepNode   `json:"direct"`
	Branches    []*stepNode `json:"branches"`
}

type expNode struct {
	Time     []wordNode `json:"time"`
	Role     []wordNode `json:"role"`
	Place    []wordNode `json:"place"`
	Input    string     `json:"input"`
	Output   string     `json:"output"`
	Material []string   `json:"material"`
	Process  *stepNode  `json:"process"`
	Name     string     `json:"name"`
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
