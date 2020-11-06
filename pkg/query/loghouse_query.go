package query

import (
	"database/sql"
	"fmt"
	qb "github.com/didi/gendry/builder"
	"github.com/sencoder/go-loghouse/pkg/uuid"
	"time"
)

const DefaultPerPage = 249

var (
	DefaultTimeParams = map[string]string{
		"format":  "seek_to",
		"seek_to": "now",
		"from":    "now-15m",
		"to":      "now",
	}
)

/*
时间参数方面，尽量简化设计
 */
type TimeParams struct {
	TimeFormat string `schema:"time_format"`
	SeekTo     string `schema:"seek_to"` // `schema:"name,required"` custom name, must be supplied
	TimeFrom   string `schema:"time_from"`
	TimeTo     string `schema:"time_to"`
}

type QueryParam struct {
	TimeParams
	Query      string   `schema:"query"`
	QueryId    string   `schema:"query_id"`
	PerPage    uint32   `schema:"per_page"`
	Namespaces []string `schema:"namespaces"`
}

type loghouseQuery struct {
	Id             string
	namespaces     []string
	from, to       time.Time
	seekTo         time.Time
	LogQueryResult LogQueryResult
	Error          error
}

func NewLoghouseQuery(query string, namespaces []string) *loghouseQuery {
	return &loghouseQuery{
		Id: uuid.TimeUUID().String(),
		LogQueryResult: LogQueryResult{
			LogEntry: make([]LogEntry, 0, DefaultPerPage),
		},
	}
}

type LogQueryResult struct {
	LogEntry []LogEntry
}

func (q *loghouseQuery) TimeParams(params TimeParams) *loghouseQuery {

	return q
}

func (q *loghouseQuery) Validate() error {
	// maybe Bad query format
	// maybe Bad time format
	return nil
}

func (q *loghouseQuery) Paginate(newerThan, olderThan string, perPage uint32) {

}

// 查询入口
func (q *loghouseQuery) Result(db sql.DB, param QueryParam) {
	where := toClickhouseWhere(
		toClickhouseParts(q.from, q.to),
		toClickhouseTimeCompare(q.from, ">"),
		toClickhouseTimeCompare(q.to, "<"),
		toClickhouseNamespaces(q.namespaces),
	)
	q.Error = q.Run(db, "logs", where, []string{"*"})
}

//
//func (q *loghouseQuery) ResultNewer() {
//
//}
//
//func (q *loghouseQuery) ResultOlder() {
//
//}
//func (q *loghouseQuery) ResultFromSeek() {
//
//}

func (q *loghouseQuery) ParseTimeFrom(from string) *time.Time {
	// set from
	return nil
}

func (q *loghouseQuery) ParseTimeTo(to string) *time.Time {
	// set to
	return nil
}

func (q *loghouseQuery) ParseTimeSeekTo(seekTo string) *time.Time {
	// set seekTo
	return nil
}

func (q *loghouseQuery) ParseTime(timeStr string) *time.Time {
	return nil
}

func (q *loghouseQuery) Where() {

}

func toClickhouseTime(t time.Time) string {
	return fmt.Sprintf("toDateTime('%s')", t.In(time.UTC).Format("2006-01-02 15:04:05"))
}

func toClickhouseParts(from, to time.Time) map[string]interface{} {

	y, m, d := to.Date()
	end := time.Date(y, m, d, 23, 59, 59, 999999999, to.Location())

	dates := make([]string, 0)
	for d := from; d.Before(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}

	return map[string]interface{}{
		"date in": dates,
	}
}

func toClickhouseTimeCompare(t time.Time, compare string) map[string]interface{} {
	if t.Nanosecond() == 0 {
		//return "timestamp " + compare + " " + toClickhouseTime(t)
		return map[string]interface{}{
			"timestamp " + compare: toClickhouseTime(t),
		}
	}

	return map[string]interface{}{
		"_or": []map[string]interface{}{
			{
				"timestamp = ":    toClickhouseTime(t),
				"nsec " + compare: t.Nanosecond(),
			},
		},
	}
}

func toClickhouseNamespaces(namespaces []string) map[string]interface{} {
	pairs := make([]map[string]interface{}, len(namespaces))
	for _, v := range namespaces {
		pairs = append(pairs, map[string]interface{}{
			"namespace =": v,
		})
	}

	return map[string]interface{}{
		"_or": pairs,
	}
}

func toClickhouseWhere(conditions ...map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{})
	for _, condition := range conditions {
		for k, v := range condition {
			ret[k] = v
		}
	}

	ret["_orderby"] = "timestamp DESC, nesc DESC"
	return ret
}

/*
where := map[string]interface{}{
    "age >": 100,
	"bar <=": 45,
    "_or": []map[string]interface{}{
        {
            "x1":    11,
            "x2 >=": 45,
        },
        {
            "x3":    "234",
            "x4 <>": "tx2",
        },
    },
    "_orderby": "fieldName asc",
    "_groupby": "fieldName",
    "_having": map[string]interface{}{"foo":"bar",},
    "_limit": []uint{offset, row_count},
    "_lockMode": "share",
}
*/

func (q *loghouseQuery) RunSql(sql string) error {
	return nil
}

func (q *loghouseQuery) Run(db sql.DB, table string, where map[string]interface{}, selectField []string) error {

	cond, vals, err := qb.BuildSelect(table, where, selectField)
	if nil != err {
		return err
	}

	rows, err := db.Query(cond, vals...)
	if nil != err {
		return err
	}
	defer rows.Close()

	for index := 0; rows.Next(); index++ {
		var entry LogEntry
		if err := rows.Scan(&entry.Labels, &entry.NumberFields, &entry.StringFields); err != nil {
			return err
			break
		}
		q.LogQueryResult.LogEntry = append(q.LogQueryResult.LogEntry, entry)
	}

	return nil
}
