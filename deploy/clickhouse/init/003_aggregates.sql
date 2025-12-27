CREATE TABLE streaming.ints_sums
(
    id UInt8,
    pos_sum Int64,
    neg_sum Int64
)
ENGINE = SummingMergeTree
ORDER BY id;
