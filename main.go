package main

import (
	"net/http"
)

func main() {
	// 第一引数 "ホスト名:ポート番号" ホスト名を省略するとlocalhostになる
	// 第二引数 HTTPハンドラを指定する。nilの場合はDefaultServeMuxが使われる
	http.ListenAndServe(":8080", nil)
}
