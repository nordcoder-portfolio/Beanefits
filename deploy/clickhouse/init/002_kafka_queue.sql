CREATE TABLE IF NOT EXISTS streaming.ints_queue
(
    raw String
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list    = 'redpanda:9092',
    kafka_topic_list     = 'ints',
    kafka_group_name     = 'ch_ints_v2',
    kafka_format         = 'LineAsString',
    kafka_num_consumers  = 1;
