package main

import (
	"flag"
	"github.com/chat/trace"
	"log"
	"net/http"
	"os"
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

// ServeHttp メソッド呼び出しする時にonceの値は常に同じものを使う必要があるので、
// このレシーバはポインタである必要がある
func (t *templateHandler) ServeHTTP(write http.ResponseWriter, request *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	err := t.templ.Execute(write, request)
	if err != nil {
		log.Fatal("template execute:", err)
	}
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	// フラグを解釈する。
	// コマンドラインで指定された文字列から必要な情報を取り出し*addrにセット。
	flag.Parse()

	r := newRoom()
	// 引数に出力先を渡す。Stdoutつまり標準出力。
	r.tracer = trace.New(os.Stdout)

	// ルート
	// *authHandlerのServeHTTPが実行され、認証が成功した場合にのみ*templateHandlerのServeHTTPが実行される
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", MustAuth(&templateHandler{filename: "login.html"}))
	http.Handle("/room", r)
	// チャットルームを開始します。
	// チャットルームはgoroutineとして実行され、チャット関連の処理はバックグラウンドで行われる。
	go r.run()
	// Webサーバを開始します
	log.Println("Webサーバを開始します。ポート: ", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
