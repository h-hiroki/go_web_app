package main

import (
	"fmt"
	"net/http"
)

func main() {
	// ルートにアクセスしたときに Hello Worldを返却する
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello World")
	})

	// 第一引数 "ホスト名:ポート番号" ホスト名を省略するとlocalhostになる
	// 第二引数 HTTPハンドラを指定する。nilの場合はDefaultServeMuxが使われる
	http.ListenAndServe(":8080", nil)
}
