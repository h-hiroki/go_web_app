package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// ResponseObject はお試し用の構造体。なにもいみはない
type ResponseObject struct {
	Status int    `json:"status"`
	Result []Task `json:"result"`
}

type Task struct {
	ID        int    `json:"id"`
	Message   string `json:"message"`
	Status    int    `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET": // index

		// DBに接続する
		// 第二引数 user:password@tcp(host:port)/dbname
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		var (
			id        int
			message   string
			status    int
			createdAt string
			updatedAt string
		)
		rows, err := db.Query(`
			SELECT
				id, message, status, created_at, updated_at
			FROM tasks
			WHERE deleted_at IS NULL
			ORDER BY id DESC
			LIMIT 20
		`)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		Tasks := []Task{}

		for rows.Next() {
			err = rows.Scan(&id, &message, &status, &createdAt, &updatedAt)
			if err != nil {
				log.Fatal(err)
			}

			Tasks = append(Tasks, Task{id, message, status, createdAt, updatedAt})
		}

		// fmt.Println(Tasks)

		getTask := ResponseObject{http.StatusOK, Tasks}

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
