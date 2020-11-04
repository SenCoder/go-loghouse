package handler

import (
	"github.com/sencoder/go-loghouse/pkg/log"
	"github.com/sencoder/go-loghouse/pkg/render"
	"github.com/sirupsen/logrus"
	"net/http"
)

//# 获取日志查询主页面
// _result_entry
// _result
// _seek_to_anchor
// _super_date_picker
// index
// layout
func QueryHandler(w http.ResponseWriter, r *http.Request) {
	queryId := r.FormValue("query_id")
	logrus.Infof("path=%s query=%v", r.URL.Path, queryId)
	// request for given query
	if queryId != "" {
		return
	}

	type TabQuery struct {
		QueryId   string
		QueryUser string
	}

	data := struct {
		Tab          string
		Username     string
		Version      string
		TabQueries   []TabQuery
		ErrorMessage string
		ParsedSeekTo string
	}{
		Tab:      "",
		Username: "admin",
		Version:  "e2282b63",
	}

	err := render.App.HTML(w, http.StatusOK, "index", data)
	if err != nil {
		log.Infof("render error: %s", err)
		//render.App.HTML(w, http.StatusOK, "error", "gophers")
	}
}
