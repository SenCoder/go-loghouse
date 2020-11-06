package clickhouse

import (
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/pkg/errors"
	"github.com/sencoder/go-loghouse/pkg/log"
)

type Config struct {
	Address                string `json:"address" yaml:"address"`
	Username               string `json:"username" yaml:"username"`
	Password               string `json:"password" yaml:"password"`
	Database               string `json:"database" yaml:"database"`
	ReadTimeout            uint64 `json:"read_timeout" yaml:"read_timeout"`   // in seconds
	WriteTimeout           uint64 `json:"write_timeout" yaml:"write_timeout"` // in seconds
	AltHosts               string `json:"alt_hosts" yaml:"alt_hosts"`
	Debug                  bool   `json:"debug" yaml:"debug"`
	ConnectionOpenStrategy string `json:"connection_open_strategy" yaml:"connection_open_strategy"`
}

/*
DSN
====================================
username/password - auth credentials
database - select the current default database
read_timeout/write_timeout - timeout in second
no_delay - disable/enable the Nagle Algorithm for tcp socket (default is 'true' - disable)
alt_hosts - comma separated list of single address host for load-balancing
connection_open_strategy - random/in_order (default random).
	random - choose random server from set
	in_order - first live server is choosen in specified order
	time_random - choose random(based on current time) server from set. This option differs from random in that randomness is based on current time rather than on amount of previous connections.
block_size - maximum rows in block (default is 1000000). If the rows are larger then the data will be split into several blocks to send them to the server. If one block was sent to the server, the data will be persisted on the server disk, we can't rollback the transaction. So always keep in mind that the batch size no larger than the block_size if you want atomic batch insert.
pool_size - maximum amount of preallocated byte chunks used in queries (default is 100). Decrease this if you experience memory problems at the expense of more GC pressure and vice versa.
debug - enable debug output (boolean value)
*/

func connect(cfg *Config) (*sql.DB, error) {

	dsn := fmt.Sprintf("tcp://%s?username=%s&password=%s&database=%s&debug=%v", cfg.Address, cfg.Username, cfg.Password, cfg.Database, cfg.Debug)
	conn, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "connect to db fail")
	}
	// todo: clickhouse-go 尚未提供自定义 Logger 接口

	if err = conn.Ping(); err != nil {
		handleException(err)
		return nil, err
	}

	return conn, nil
}

func disconnect(conn *sql.DB) {
	if conn != nil {
		conn.Close()
	}
}

type Initializer struct {
	*sql.DB
	database string
	cluster  string
	ttl      uint
}

func NewInitializer(cfg *Config, cluster string, ttl uint) *Initializer {

	var db = DBName
	switch {
	case cfg != nil:
		if cfg.Database != "" {
			db = cfg.Database
		}
	}

	return &Initializer{
		database: db,
		cluster:  cluster,
		ttl:      ttl,
	}
}

func (l *Initializer) CreateMergeTreeTable() {
	sqlStr := fmt.Sprintf(logsRpl, l.cluster, l.database)
	_, err := l.Exec(sqlStr)
	handleCreateException(err)
}

func (l *Initializer) CreateBufferTable() {
	sqlStr := fmt.Sprintf(logsBuffer, l.cluster, l.database)
	_, err := l.Exec(sqlStr)
	handleCreateException(err)
}

func (l *Initializer) CreateDistributedTable() {
	sqlStr := fmt.Sprintf(logsD, l.cluster, l.database, l.cluster)
	_, err := l.Exec(sqlStr)
	handleCreateException(err)
}

func handleCreateException(err error) {
	if err == nil {
		return
	}
	if exception, ok := err.(*clickhouse.Exception); ok {
		log.Fatalf("Code: %d. Exception: %s. Stack: %s", exception.Code, exception.Message, exception.StackTrace)
	}
}

func handleException(err error) {
	if err == nil {
		return
	}
	if exception, ok := err.(*clickhouse.Exception); ok {
		log.Warnf("Code: %d. Exception: %s. Stack: %s", exception.Code, exception.Message, exception.StackTrace)
	}
}
