package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func serveMain(serveString string) {
	//http.HandleFunc("/explain", explainHttpHandler)
	http.HandleFunc("/analyse", analyseHttpHandler)
	http.HandleFunc("/extract", extractHttpHandler)

	fmt.Println("Serving at " + serveString)

	err := http.ListenAndServe(serveString, nil)
	if err != nil {
		fmt.Println("error:", err)
	}
}

//func explainHttpHandler(w http.ResponseWriter, r *http.Request) {
//	// Cross-domain
//	w.Header().Set("Access-Control-Allow-Origin", "*")
//	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
//	w.Header().Set("content-type", "application/json")
//
//	if r.Method == "POST" {
//		var cm clientMessage
//		decoder := json.NewDecoder(r.Body)
//
//		err := decoder.Decode(&cm)
//		if err != nil {
//			sr := serverRespond{
//				Code:   500,
//				Status: false,
//				Method: r.Method,
//				Text:   "JSON parse error: " + err.Error(),
//			}
//			stream, err := json.Marshal(sr)
//			if err != nil {
//				fmt.Println("error:", err)
//			}
//			_, _ = fmt.Fprintf(w, string(stream))
//		}
//
//		// Deal with panic(s)
//		defer func() {
//			if err := recover(); err != nil {
//				errorContent, flag := err.(string)
//				if flag {
//					sr := serverRespond{
//						Code:   500,
//						Status: false,
//						Method: r.Method,
//						Text:   errorContent,
//					}
//					stream, err := json.Marshal(sr)
//					if err != nil {
//						fmt.Println("error:", err)
//					}
//					_, _ = fmt.Fprintf(w, string(stream))
//				}
//			}
//		}()
//
//		group := extract(cm.Content)
//		node, exp := analyse(group)
//
//		if cm.Target == "" {
//			panic("missing target to explain")
//		} else {
//			out := string(explain(node, exp, cm.Target))
//			_, _ = fmt.Fprintf(w, out)
//		}
//	} else {
//		methodNotAllow(w, r)
//	}
//}
//
func analyseHttpHandler(w http.ResponseWriter, r *http.Request) {
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
			_, _ = fmt.Fprintf(w, string(stream))
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
					_, _ = fmt.Fprintf(w, string(stream))
				}
			}
		}()

		var sr serverRespond
		group, pErr := extract(cm.Content)
		if pErr != nil {
			out := stringifyParseError(pErr)
			sr = serverRespond{
				Code:   200,
				Status: false,
				Method: r.Method,
				Text:   out,
			}
		} else {
			nodes, pErr := analyse(group)
			if pErr != nil {
				out := stringifyParseError(pErr)
				sr = serverRespond{
					Code:   200,
					Status: false,
					Method: r.Method,
					Text:   out,
				}
			} else {
				out := stringifyExpNodes(nodes)
				sr = serverRespond{
					Code:   200,
					Status: true,
					Method: r.Method,
					Text:   out,
				}
			}
		}

		stream, err := json.Marshal(sr)
		if err != nil {
			fmt.Println("error:", err)
		}
		_, _ = fmt.Fprintf(w, string(stream))
	} else {
		methodNotAllow(w, r)
	}
}

func extractHttpHandler(w http.ResponseWriter, r *http.Request) {
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
			_, _ = fmt.Fprintf(w, string(stream))
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
					_, _ = fmt.Fprintf(w, string(stream))
				}
			}
		}()

		var sr serverRespond
		group, pErr := extract(cm.Content)
		if pErr != nil {
			out := stringifyParseError(pErr)
			sr = serverRespond{
				Code:   200,
				Status: false,
				Method: r.Method,
				Text:   out,
			}
		} else {
			out := stringifyWordNodes(group)
			sr = serverRespond{
				Code:   200,
				Status: true,
				Method: r.Method,
				Text:   out,
			}
		}

		stream, err := json.Marshal(sr)
		if err != nil {
			fmt.Println("error:", err)
		}
		_, _ = fmt.Fprintf(w, string(stream))
	} else {
		methodNotAllow(w, r)
	}
}
