package main

import (
	"fmt"
	"github.com/stretchr/gomniauth"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// authというkeyでクッキーを取り出す
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 未認証
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// 何らかの別のエラーが発生
		panic(err.Error())
	} else {
		// 成功。ラップされたハンドラを呼び出します
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler はサードパーティへのログインの処理を待ち受ける。
// 内部状態を保持する必要がない為、http.Handlerインタフェースは実装していない。
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// パスの文字列を分割
	// パスの形式: /auth/{action}/{provider}
	segs := strings.Split(r.URL.Path, "/")
	// FIXME segs[2]とsegs[3]が必ず存在すると仮定している為、不完全なパスが指定されているとコネクションが切断される
	// 例えば、/auth/nonsense のようなアクセスだと切れる。
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました:", provider, "-", err)
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("GetBeginAuthURLの呼び出し中にエラーが発生しました:", provider, "-",err)
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}
