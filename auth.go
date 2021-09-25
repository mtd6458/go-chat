package main

import "net/http"

type authHandler struct {
  next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
  if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
    // 未認証
    w.Header().Set("Location", "/login")
    w.WriteHeader(http.StatusTemporaryRedirect)
  }
}
