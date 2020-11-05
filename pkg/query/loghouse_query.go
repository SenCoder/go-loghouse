package query

const DefaultPerPage = 249

var (
	DefaultTimeParams = map[string]string{
		"format":  "seek_to",
		"seek_to": "now",
		"from":    "now-15m",
		"to":      "now",
	}
)

type TimeParams struct {
	TimeFormat string   `schema:"time_format"`
	SeekTo     string   `schema:"seek_to"` // `schema:"name,required"` custom name, must be supplied
	TimeFrom   string   `schema:"time_from"`
	TimeTo     string   `schema:"time_to"`
}

type loghouseQuery struct {
	Id string
	Namespaces []string
	OrderBy string
}

func NewLoghouseQuery(query string, namespaces []string) *loghouseQuery {
	return &loghouseQuery{

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

func (q *loghouseQuery) Result() {

}