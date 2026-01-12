CREATE MATERIALIZED VIEW IF NOT EXISTS streaming.mv_ints_to_sums
            TO streaming.ints_sums
AS
SELECT
    1 AS id,
    if(n > 0, n, 0) AS pos_sum,
    if(n < 0, n, 0) AS neg_sum
FROM
    (
        SELECT toInt64OrNull(trim(raw)) AS n
        FROM streaming.ints_queue
    )
WHERE n IS NOT NULL;


CREATE MATERIALIZED VIEW IF NOT EXISTS streaming.mv_ints_to_dlq
            TO streaming.ints_dlq
AS
SELECT
    now64(3) AS event_time,
    _topic AS topic,
    toUInt64(_partition) AS partition,
    toUInt64(_offset) AS offset,
    ifNull(_key, '') AS key,
    raw AS raw_message,
    'invalid_int64' AS error
FROM streaming.ints_queue
WHERE toInt64OrNull(trim(raw)) IS NULL;
