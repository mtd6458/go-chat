package main

import (
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
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
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました:", p, "-", err)
		}
		loginUrl, err := p.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("GetBeginAuthURLの呼び出し中にエラーが発生しました:", p, "-",err)
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		p, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました:", p, "-", err)
		}
		creds, err := p.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			log.Fatalln("認証を完了出来ませんでした:", p, "-", err)
		}

		user, err := p.GetUser(creds)
		if err != nil {
			log.Fatalln("ユーザの取得に失敗しました:", p, "-", err)

		}
		authCookeieValue := objx.New(map[string]interface{} {
			"name": user.Name(),
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:       "auth",
			Value:      authCookeieValue,
			Path:       "/",
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}
