package handler

import (
	"github.com/gorilla/schema"
	"github.com/sencoder/go-loghouse/pkg/log"
	"github.com/sencoder/go-loghouse/pkg/query"
	"github.com/sencoder/go-loghouse/pkg/render"
	"github.com/sirupsen/logrus"
	"net/http"
)

var decoder = schema.NewDecoder()

var datePicker = SuperDatepickerPresets{
	LastDays: []Last{
		{Name: "Last 2 days", From: "now-2d", To: "now"},
		{Name: "Last 3 days", From: "now-3d", To: "now"},
		{Name: "Last 7 days", From: "now-7d", To: "now"},
		{Name: "Last 14 days", From: "now-14d", To: "now"},
		{Name: "Last 30 days", From: "now-30d", To: "now"},
		{Name: "Last 60 days", From: "now-60d", To: "now"},
		{Name: "Last 90 days", From: "now-90d", To: "now"},
	},
	LastDay: []Last{
		{Name: "Last 5 min", From: "now-5m", To: "now"},
		{Name: "Last 15 min", From: "now-15m", To: "now"},
		{Name: "Last 30 min", From: "now-30m", To: "now"},
		{Name: "Last 1 hour", From: "now-1h", To: "now"},
		{Name: "Last 3 hour", From: "now-3h", To: "now"},
		{Name: "Last 6 hour", From: "now-6h", To: "now"},
		{Name: "Last 12 hour", From: "now-12h", To: "now"},
		{Name: "Last 24 hour", From: "now-24h", To: "now"},
	},
	SeekTo: []string{
		"5 minutes ago",
		"15 minutes ago",
		"30 minutes ago",
		"1 hour ago",
		"3 hours ago",
		"6 hours ago",
		"12 hours ago",
		"24 hours ago",
	},
}



type Last struct {
	Name string
	From string
	To   string
}

type SuperDatepickerPresets struct {
	LastDays []Last
	LastDay  []Last
	SeekTo   []string
}

type QueryParam struct {
	query.TimeParams
	Query      string   `schema:"query"`
	QueryId    string   `schema:"query_id"`
	PerPage    string   `schema:"per_page"`
	Namespaces []string `schema:"namespaces"`
}


//# 获取日志查询主页面
// _result_entry
// _result
// _seek_to_anchor
// _super_date_picker
// index
// layout
func QueryHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		// Handle error
	}

	var queryParam QueryParam

	logrus.Infof("path=%s query=%+v", r.URL.Path, queryParam)
	if err := decoder.Decode(&queryParam, r.Form); err != nil {
		// Handle error
	}

	logrus.Infof("path=%s query=%+v", r.URL.Path, queryParam)
	// request for given query
	if queryParam.QueryId != "" {
		// query by id
		return
	}

	// query from params

	type TabQuery struct {
		QueryId   string
		QueryUser string
	}

	type Attribute struct {
		Query string
	}

	data := struct {
		Tab                 string
		Username            string
		AvailableNamespaces []string
		Version             string
		TabQueries          []TabQuery
		ErrorMessage        string
		ParsedSeekTo        string
		Result              []query.LogQueryResult
		Attribute           Attribute
		Namespaces          []string
		Persisted           bool
		Params              map[string]string
		TimeParams          map[string]string
		SuperDatepickerPresets
	}{
		Tab:                    "",
		Username:               "admin",
		AvailableNamespaces:    []string{"review-123", "production", "staging"},
		Version:                "e2282b63",
		Namespaces:             []string{"staging"},
		Params:                 map[string]string{"format": "rangy"},
		TimeParams:             query.DefaultTimeParams,
		SuperDatepickerPresets: datePicker,
	}

	for k := range data.TimeParams {
		switch k {
		case "format":
			data.TimeParams[k] = queryParam.TimeFormat
		case "seek_to":
			data.TimeParams[k] = queryParam.SeekTo
		case "from":
			data.TimeParams[k] = queryParam.TimeFrom
		case "to":
			data.TimeParams[k] = queryParam.TimeTo
		}
	}

	q := query.NewLoghouseQuery(queryParam.Query, queryParam.Namespaces).TimeParams(queryParam.TimeParams)
	if err := q.Validate(); err != nil {

	}


	err := render.App.HTML(w, http.StatusOK, "index", data)
	if err != nil {
		log.Infof("render error: %s", err)
		//render.App.HTML(w, http.StatusOK, "error", "gophers")
	}
}
