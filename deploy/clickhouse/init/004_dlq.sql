CREATE TABLE streaming.ints_dlq
(
    event_time DateTime64(3),
    topic LowCardinality(String),
    partition UInt64,
    offset UInt64,
    key String,
    raw_message String,
    error String
)
ENGINE = MergeTree
ORDER BY (event_time, topic, partition, offset);
