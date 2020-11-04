package handler

import (
	"net/http"
)

//# 对路由 / 的 get 请求被重定向到 /query
func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/query", http.StatusFound)
}
