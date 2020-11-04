package handler

import (
	"github.com/sencoder/go-loghouse/httpd/response"
	"github.com/sencoder/go-loghouse/pkg/render"
	"github.com/sencoder/go-loghouse/version"
	"net/http"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {

	data := struct {
		BuildTime    string `json:"build_time"`
		BuildVersion string `json:"build_version"`
	}{
		BuildTime:    version.BuildTime,
		BuildVersion: version.BuildVersion,
	}

	_ = render.App.JSON(w, http.StatusOK, response.Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}
