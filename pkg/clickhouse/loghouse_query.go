package clickhouse

type LoghouseQuery struct {

}

type QueryResult struct {
	LogEntry []LogEntry
}

func (q *LoghouseQuery) Result() {

}