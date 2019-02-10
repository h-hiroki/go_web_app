package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// ResponseObject レスポンス共通の形式を構築する
type ResponseObject struct {
	Status int    `json:"status"`
	Result []Task `json:"result"`
}

// Task タスクデータを構築する
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
		var (
			id        int
			message   string
			status    int
			createdAt string
			updatedAt string
		)
		for rows.Next() {
			err = rows.Scan(&id, &message, &status, &createdAt, &updatedAt)
			if err != nil {
				log.Fatal(err)
			}

			Tasks = append(Tasks, Task{id, message, status, createdAt, updatedAt})
		}

		getTask := ResponseObject{http.StatusOK, Tasks}

		res, err := json.Marshal(getTask)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	case "POST": // create
		// DBに接続する
		// 第二引数 user:password@tcp(host:port)/dbname
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// insertする
		// TODO requestのbodyを読み込んでそのメッセージを更新する
		_, err = db.Exec(`
		INSERT INTO tasks(message, status, created_at, updated_at) VALUES('WebAPI inserted.', 0, now(), now())
		`)
		if err != nil {
			log.Fatal(err)
		}

		// 最新のtasksテーブルの結果を取得する
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
		var (
			id        int
			message   string
			status    int
			createdAt string
			updatedAt string
		)
		for rows.Next() {
			err = rows.Scan(&id, &message, &status, &createdAt, &updatedAt)
			if err != nil {
				log.Fatal(err)
			}

			Tasks = append(Tasks, Task{id, message, status, createdAt, updatedAt})
		}

		getTask := ResponseObject{http.StatusOK, Tasks}

		res, err := json.Marshal(getTask)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func taskResorceHandler(w http.ResponseWriter, r *http.Request) {
	// URLパラメタを取得する
	URLPathParts := strings.Split(r.URL.Path, "/")
	TaskID := URLPathParts[len(URLPathParts)-1]

	switch r.Method {
	case "GET": // show
		// DBに接続する
		// 第二引数 user:password@tcp(host:port)/dbname
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// 最新のtasksテーブルの結果を取得する
		query := "SELECT id, message, status, created_at, updated_at FROM tasks WHERE deleted_at IS NULL AND id = " + TaskID

		rows, err := db.Query(query)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		Tasks := []Task{}
		var (
			id        int
			message   string
			status    int
			createdAt string
			updatedAt string
		)
		for rows.Next() {
			err = rows.Scan(&id, &message, &status, &createdAt, &updatedAt)
			if err != nil {
				log.Fatal(err)
			}

			Tasks = append(Tasks, Task{id, message, status, createdAt, updatedAt})
		}

		getTask := ResponseObject{http.StatusOK, Tasks}

		res, err := json.Marshal(getTask)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	case "POST": // update
		fmt.Println("1件のTODOを更新")
	case "DELETE": // destroy
		// DBに接続する
		// 第二引数 user:password@tcp(host:port)/dbname
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Deleteする
		query := "UPDATE tasks set deleted_at = now() where id = " + TaskID
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal(err)
		}

		// 最新のtasksテーブルの結果を取得する
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
		var (
			id        int
			message   string
			status    int
			createdAt string
			updatedAt string
		)
		for rows.Next() {
			err = rows.Scan(&id, &message, &status, &createdAt, &updatedAt)
			if err != nil {
				log.Fatal(err)
			}

			Tasks = append(Tasks, Task{id, message, status, createdAt, updatedAt})
		}

		getTask := ResponseObject{http.StatusOK, Tasks}

		res, err := json.Marshal(getTask)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
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
