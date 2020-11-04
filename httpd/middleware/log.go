/*
 * Copyright (c) 2018 PingAn. All rights reserved.
 */

package middleware

import (
	"github.com/sencoder/go-loghouse/pkg/log"
	"net/http"
)


func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Infof("remote_addr=%s method=%s path=%s", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)

	})
}