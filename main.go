package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// templ は1つのテンプレートを表します
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHttp メソッド呼び出しする時にonceの値は常に同じものを使う必要があるので、このレシーバはポインタである必要がある
func (t *templateHandler) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	err := t.templ.Execute(write, nil)
	if err != nil {
		log.Fatal("template execute:", err)
	}
}

func main() {
	// ルート
	http.Handle("/", &templateHandler{filename: "chat.html"})
	// Webサーバを開始します
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
