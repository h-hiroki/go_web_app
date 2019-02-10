package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ResponseObject はお試し用の構造体。なにもいみはない
type ResponseObject struct {
	Status int    `json:"status"`
	Result string `json:"result"`
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // index
		getTask := ResponseObject{http.StatusOK, "OK"}

		res, err := json.Marshal(getTask)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	case "POST": // create
		fmt.Println("新規TODO作成")
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func taskResorceHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // show
		fmt.Println("1件のTODO詳細を返却")
	case "POST": // update
		fmt.Println("1件のTODOを更新")
	case "DELETE": // destroy
		fmt.Println("1件のTODOを削除")
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	// TODO RESTfulで作りたい
	http.HandleFunc("/task", taskHandler)
	http.HandleFunc("/task/", taskResorceHandler)

	http.ListenAndServe(":8080", nil)
}
