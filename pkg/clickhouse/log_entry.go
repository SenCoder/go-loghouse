package clickhouse

type LogEntry struct {
	Labels       []string
	StringFields []string
	NumberFields []string
}
