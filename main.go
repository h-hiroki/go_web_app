package main

import (
	"encoding/json"
	"io/ioutil"
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

// CreateTask 新規TODO作成時のメッセージを受信する
type CreateTask struct {
	Message string `json:"message"`
}

// UpdateTask TODO更新時の各種パラメタを受信する
type UpdateTask struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
	Status  int    `json:"status"`
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
			LIMIT 30
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
		// read body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Unmarshal
		var createTask CreateTask
		err = json.Unmarshal(body, &createTask)
		if err != nil {
			log.Fatal(err)
		}

		// DBに接続する
		// 第二引数 user:password@tcp(host:port)/dbname
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// insertする
		_, err = db.Exec("INSERT INTO tasks(message, status, created_at, updated_at) VALUES(?, 0, now(), now())", createTask.Message)
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
			LIMIT 30
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
		// read body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Unmarshal
		var updateTask UpdateTask
		err = json.Unmarshal(body, &updateTask)
		if err != nil {
			log.Fatal(err)
		}

		// DBに接続する
		// 第二引数 user:password@tcp(host:port)/dbname
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_db")
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// insertする
		_, err = db.Exec("UPDATE tasks SET message = ?, status = ?, updated_at = now() where id = ?", updateTask.Message, updateTask.Status, updateTask.ID)
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
			LIMIT 30
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
			LIMIT 30
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
