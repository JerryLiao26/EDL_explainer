package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func serveMain(serveString string) {
	http.HandleFunc("/explain", explainHandler)
	http.HandleFunc("/analyse", analyseHandler)

	fmt.Println("Serving at " + serveString)
	err := http.ListenAndServe(serveString, nil)
	if err != nil {
		fmt.Println("error:", err)
	}
}

func explainHandler(w http.ResponseWriter, r *http.Request) {
	// Cross-domain
	w.Header().Set("Access-Control-Allow-Origin", "*")             // 允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") // header的类型
	w.Header().Set("content-type", "application/json")             // 返回数据格式是json
	if r.Method == "POST" {
		var cm clientMessage
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&cm)
		if err != nil {
			sr := serverRespond{
				Code:   500,
				Status: false,
				Method: r.Method,
				Text:   "JSON parse error: " + err.Error(),
			}
			stream, err := json.Marshal(sr)
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Fprintf(w, string(stream))
		}
		// Deal with panic(s)
		defer func() {
			if err := recover(); err != nil {
				errorContent, flag := err.(string)
				if flag {
					sr := serverRespond{
						Code:   500,
						Status: false,
						Method: r.Method,
						Text:   errorContent,
					}
					stream, err := json.Marshal(sr)
					if err != nil {
						fmt.Println("error:", err)
					}
					fmt.Fprintf(w, string(stream))
				}
			}
		}()
		group := extract(cm.Content)
		node, exp := analyse(group)
		if cm.Target == "" {
			panic("missing target to explain")
		} else {
			out := string(explain(node, exp, cm.Target))
			fmt.Fprintf(w, out)
		}
	} else {
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
		fmt.Fprintf(w, string(stream))
	}
}

func analyseHandler(w http.ResponseWriter, r *http.Request) {
	// Cross-domain
	w.Header().Set("Access-Control-Allow-Origin", "*")             // 允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") // header的类型
	w.Header().Set("content-type", "application/json")             // 返回数据格式是json
	if r.Method == "POST" {
		var cm clientMessage
		decoder := json.NewDecoder(r.Body)
		parseErr := decoder.Decode(&cm)
		if parseErr != nil {
			sr := serverRespond{
				Code:   500,
				Status: false,
				Method: r.Method,
				Text:   "JSON parse error: " + parseErr.Error(),
			}
			stream, err := json.Marshal(sr)
			if err != nil {
				fmt.Println("error:", err)
			}
			fmt.Fprintf(w, string(stream))
		}
		// Deal with panic(s)
		defer func() {
			if err := recover(); err != nil {
				errorContent, flag := err.(string)
				if flag {
					sr := serverRespond{
						Code:   500,
						Status: false,
						Method: r.Method,
						Text:   errorContent,
					}
					stream, err := json.Marshal(sr)
					if err != nil {
						fmt.Println("error:", err)
					}
					fmt.Fprintf(w, string(stream))
				}
			}
		}()
		group := extract(cm.Content)
		node, exp := analyse(group)
		out := "[" + string(node) + "," + string(exp) + "]"
		sr := serverRespond{
			Code:   200,
			Status: true,
			Method: r.Method,
			Text:   out,
		}
		stream, err := json.Marshal(sr)
		if err != nil {
			fmt.Println("error:", err)
		}
		fmt.Fprintf(w, string(stream))
	} else {
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
		fmt.Fprintf(w, string(stream))
	}
}
