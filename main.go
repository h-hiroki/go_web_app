package main

import (
	"fmt"
	"net/http"
)

func paidHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "You are paid\n")
}

func notPayHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusPaymentRequired) // httpには各HTTPステータスが定数で定義されている！
	fmt.Fprint(w, "👯‍♀️Payment Required👯‍♀️\n")
}

func main() {
	// /paidにアクセスした場合の処理
	http.HandleFunc("/paid", paidHandler)
	// /not_payにアクセスした場合の処理
	http.HandleFunc("/not_pay", notPayHandler)

	// 第一引数 "ホスト名:ポート番号" ホスト名を省略するとlocalhostになる
	// 第二引数 HTTPハンドラを指定する。nilの場合はDefaultServeMuxが使われる
	http.ListenAndServe(":8080", nil)
}
