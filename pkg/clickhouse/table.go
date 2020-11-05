package clickhouse

const DBName = "logs"

// todo: 表参数 根据生产环境实际情况调整设置
// TTL 设置数据过期时间
// storage_policy 设置存储策略

const (
	logsRpl = `
CREATE TABLE IF NOT EXISTS logs_rpl
ON CLUSTER %s
(
    date 					Date MATERIALIZED toDate(timestamp),
    timestamp 				DateTime, 
    nsec 					UInt32, 
    source 					String, 
    namespace 				String, 
    host 					String, 
    pod_name 				String, 
    container_name 			String, 
    stream 					String, 
    labels.names 			Array(String), 
    labels.values 			Array(String), 
    string_fields.names 	Array(String), 
    string_fields.values 	Array(String), 
    number_fields.names 	Array(String), 
    number_fields.values 	Array(Float64), 
    boolean_fields.names 	Array(String), 
    boolean_fields.values 	Array(Float64), 
    null_fields.names 		Array(String)
)
ENGINE = ReplicatedMergeTree('/clickhouse/tables/%s/{shard}/logs_rpl', '{replica}')
PARTITION BY (date)
ORDER BY (timestamp, nsec, namespace, container_name)
TTL date + toIntervalDay(14)
SETTINGS index_granularity = 32768;
`
	logsD = `
CREATE TABLE IF NOT EXISTS logs 
ON CLUSTER %s
AS logs_rpl
ENGINE = Distributed(%s, %s, logs_rpl, rand());
`
	logsBuffer=`
CREATE TABLE IF NOT EXISTS logs_buffer
ON CLUSTER %s
AS logs_rpl
ENGINE = Buffer(currentDatabase(), logs, 16, 10, 60, 1000, 10000, 1048576, 10485760);
`
	queries = `
CREATE TABLE queries_v1
ON CLUSTER %s
(
    id FixedString(36), 
    name String, 
    namespaces Array(String), 
    query String, 
    position UInt8
)
ENGINE = ReplicatedMergeTree('/clickhouse/tables/logs/queries_v1', '{replica}')
PARTITION BY id
ORDER BY id
SETTINGS index_granularity = 8192;
`
	migrations=`
CREATE TABLE migrations
ON CLUSTER %s
(
    timestamp DateTime, 
    version UInt32
)
ENGINE = ReplicatedMergeTree('/clickhouse/tables/logs/migrations', '{replica}')
PARTITION BY timestamp
ORDER BY timestamp
SETTINGS index_granularity = 8192;
`
)
