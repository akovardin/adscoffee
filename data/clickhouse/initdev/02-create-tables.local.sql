------ requests

CREATE TABLE IF NOT EXISTS analytics.kafka_requests_raw
(
    raw_data String
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'kafka1:9092,kafka2:9092,kafka3:9092',
    kafka_topic_list = 'request',
    kafka_group_name = 'clickhouse_request_consumer',
    kafka_format = 'JSONAsString',
    kafka_max_block_size = 131072,
    kafka_poll_max_batch_size = 2000,
    kafka_flush_interval_ms = 500,
    kafka_num_consumers = 1,
    kafka_thread_per_consumer = 1;

CREATE TABLE IF NOT EXISTS analytics.requests
(
    id String,
    action String,
    timestamp DateTime,
    click_id String,
    banner_id String,
    group_id String,
    campaign_id String,
    advertiser_id String,
    gaid String,
    oaid String,
    user_id String,
    stable_id String,
    bundle String,
    city String,
    country String,
    region String,
    price Decimal64(3),
    network String
)
ENGINE = MergeTree()
ORDER BY (timestamp, id)
TTL timestamp + INTERVAL 3 DAY DELETE;

CREATE MATERIALIZED VIEW IF NOT EXISTS analytics.requests_parsed_mv
TO analytics.requests AS
SELECT 
    JSONExtractString(raw_data, 'request_id') as id,
    parseDateTimeBestEffortOrZero(JSONExtractString(raw_data, 'timestamp')) as timestamp,
    JSONExtractString(raw_data, 'action') as action,
    JSONExtractString(raw_data, 'click_id') as click_id,
    JSONExtractString(raw_data, 'banner_id') as banner_id,
    JSONExtractString(raw_data, 'group_id') as group_id,
    JSONExtractString(raw_data, 'campaign_id') as campaign_id,
    JSONExtractString(raw_data, 'advertiser_id') as advertiser_id,
    JSONExtractString(raw_data, 'gaid') as gaid,
    JSONExtractString(raw_data, 'oaid') as oaid,
    JSONExtractString(raw_data, 'user_id') as user_id,
    JSONExtractString(raw_data, 'stable_id') as stable_id,
    JSONExtractString(raw_data, 'bundle') as bundle,
    JSONExtractString(raw_data, 'city') as city,
    JSONExtractString(raw_data, 'country') as country,
    JSONExtractString(raw_data, 'region') as region,
    if(empty(JSONExtractString(raw_data, 'price')), 0, 
       toDecimal64OrNull(JSONExtractString(raw_data, 'price'), 3)) as price,
    JSONExtractString(raw_data, 'network') as network
FROM analytics.kafka_requests_raw;

------ impressions


CREATE TABLE IF NOT EXISTS analytics.kafka_impressions_raw
(
    raw_data String
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'kafka1:9092,kafka2:9092,kafka3:9092',
    kafka_topic_list = 'impression',
    kafka_group_name = 'clickhouse_impression_consumer',
    kafka_format = 'JSONAsString',
    kafka_max_block_size = 131072,
    kafka_poll_max_batch_size = 2000,
    kafka_flush_interval_ms = 500,
    kafka_num_consumers = 1,
    kafka_thread_per_consumer = 1;

CREATE TABLE IF NOT EXISTS analytics.impressions
(
    id String,
    action String,
    timestamp DateTime,
    click_id String,
    banner_id String,
    group_id String,
    campaign_id String,
    advertiser_id String,
    gaid String,
    oaid String,
    user_id String,
    stable_id String,
    bundle String,
    city String,
    country String,
    region String,
    price Decimal64(3),
    network String
)
ENGINE = MergeTree()
ORDER BY (timestamp, id)
TTL timestamp + INTERVAL 3 DAY DELETE;

CREATE MATERIALIZED VIEW IF NOT EXISTS analytics.impressions_parsed_mv
TO analytics.impressions AS
SELECT 
    JSONExtractString(raw_data, 'request_id') as id,
    parseDateTimeBestEffortOrZero(JSONExtractString(raw_data, 'timestamp')) as timestamp,
    JSONExtractString(raw_data, 'action') as action,
    JSONExtractString(raw_data, 'click_id') as click_id,
    JSONExtractString(raw_data, 'banner_id') as banner_id,
    JSONExtractString(raw_data, 'group_id') as group_id,
    JSONExtractString(raw_data, 'campaign_id') as campaign_id,
    JSONExtractString(raw_data, 'advertiser_id') as advertiser_id,
    JSONExtractString(raw_data, 'gaid') as gaid,
    JSONExtractString(raw_data, 'oaid') as oaid,
    JSONExtractString(raw_data, 'user_id') as user_id,
    JSONExtractString(raw_data, 'stable_id') as stable_id,
    JSONExtractString(raw_data, 'bundle') as bundle,
    JSONExtractString(raw_data, 'city') as city,
    JSONExtractString(raw_data, 'country') as country,
    JSONExtractString(raw_data, 'region') as region,
    if(empty(JSONExtractString(raw_data, 'price')), 0, 
       toDecimal64OrNull(JSONExtractString(raw_data, 'price'), 3)) as price,
    JSONExtractString(raw_data, 'network') as network
FROM analytics.kafka_impressions_raw;

------ clicks


CREATE TABLE IF NOT EXISTS analytics.kafka_clicks_raw
(
    raw_data String
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'kafka1:9092,kafka2:9092,kafka3:9092',
    kafka_topic_list = 'click',
    kafka_group_name = 'clickhouse_click_consumer',
    kafka_format = 'JSONAsString',
    kafka_max_block_size = 131072,
    kafka_poll_max_batch_size = 2000,
    kafka_flush_interval_ms = 500,
    kafka_num_consumers = 1,
    kafka_thread_per_consumer = 1;

CREATE TABLE IF NOT EXISTS analytics.clicks
(
    id String,
    action String,
    timestamp DateTime,
    click_id String,
    banner_id String,
    group_id String,
    campaign_id String,
    advertiser_id String,
    gaid String,
    oaid String,
    user_id String,
    stable_id String,
    bundle String,
    city String,
    country String,
    region String,
    price Decimal64(3),
    network String
)
ENGINE = MergeTree()
ORDER BY (timestamp, id)
TTL timestamp + INTERVAL 3 DAY DELETE;

CREATE MATERIALIZED VIEW IF NOT EXISTS analytics.clicks_parsed_mv
TO analytics.clicks AS
SELECT 
    JSONExtractString(raw_data, 'request_id') as id,
    parseDateTimeBestEffortOrZero(JSONExtractString(raw_data, 'timestamp')) as timestamp,
    JSONExtractString(raw_data, 'action') as action,
    JSONExtractString(raw_data, 'click_id') as click_id,
    JSONExtractString(raw_data, 'banner_id') as banner_id,
    JSONExtractString(raw_data, 'group_id') as group_id,
    JSONExtractString(raw_data, 'campaign_id') as campaign_id,
    JSONExtractString(raw_data, 'advertiser_id') as advertiser_id,
    JSONExtractString(raw_data, 'gaid') as gaid,
    JSONExtractString(raw_data, 'oaid') as oaid,
    JSONExtractString(raw_data, 'user_id') as user_id,
    JSONExtractString(raw_data, 'stable_id') as stable_id,
    JSONExtractString(raw_data, 'bundle') as bundle,
    JSONExtractString(raw_data, 'city') as city,
    JSONExtractString(raw_data, 'country') as country,
    JSONExtractString(raw_data, 'region') as region,
    if(empty(JSONExtractString(raw_data, 'price')), 0, 
       toDecimal64OrNull(JSONExtractString(raw_data, 'price'), 3)) as price,
    JSONExtractString(raw_data, 'network') as network
FROM analytics.kafka_clicks_raw;

------ responses


CREATE TABLE IF NOT EXISTS analytics.kafka_responses_raw
(
    raw_data String
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'kafka1:9092,kafka2:9092,kafka3:9092',
    kafka_topic_list = 'response',
    kafka_group_name = 'clickhouse_response_consumer',
    kafka_format = 'JSONAsString',
    kafka_max_block_size = 131072,
    kafka_poll_max_batch_size = 2000,
    kafka_flush_interval_ms = 500,
    kafka_num_consumers = 1,
    kafka_thread_per_consumer = 1;

CREATE TABLE IF NOT EXISTS analytics.responses
(
    id String,
    action String,
    timestamp DateTime,
    click_id String,
    banner_id String,
    group_id String,
    campaign_id String,
    advertiser_id String,
    gaid String,
    oaid String,
    user_id String,
    stable_id String,
    bundle String,
    city String,
    country String,
    region String,
    price Decimal64(3),
    network String
)
ENGINE = MergeTree()
ORDER BY (timestamp, id)
TTL timestamp + INTERVAL 3 DAY DELETE;

CREATE MATERIALIZED VIEW IF NOT EXISTS analytics.responses_parsed_mv
TO analytics.responses AS
SELECT 
    JSONExtractString(raw_data, 'request_id') as id,
    parseDateTimeBestEffortOrZero(JSONExtractString(raw_data, 'timestamp')) as timestamp,
    JSONExtractString(raw_data, 'action') as action,
    JSONExtractString(raw_data, 'click_id') as click_id,
    JSONExtractString(raw_data, 'banner_id') as banner_id,
    JSONExtractString(raw_data, 'group_id') as group_id,
    JSONExtractString(raw_data, 'campaign_id') as campaign_id,
    JSONExtractString(raw_data, 'advertiser_id') as advertiser_id,
    JSONExtractString(raw_data, 'gaid') as gaid,
    JSONExtractString(raw_data, 'oaid') as oaid,
    JSONExtractString(raw_data, 'user_id') as user_id,
    JSONExtractString(raw_data, 'stable_id') as stable_id,
    JSONExtractString(raw_data, 'bundle') as bundle,
    JSONExtractString(raw_data, 'city') as city,
    JSONExtractString(raw_data, 'country') as country,
    if(empty(JSONExtractString(raw_data, 'price')), 0, 
       toDecimal64OrNull(JSONExtractString(raw_data, 'price'), 3)) as price,
    JSONExtractString(raw_data, 'network') as network
FROM analytics.kafka_responses_raw;

------ conversions


CREATE TABLE IF NOT EXISTS analytics.kafka_conversions_raw
(
    raw_data String
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'kafka1:9092,kafka2:9092,kafka3:9092',
    kafka_topic_list = 'conversion',
    kafka_group_name = 'clickhouse_conversion_consumer',
    kafka_format = 'JSONAsString',
    kafka_max_block_size = 131072,
    kafka_poll_max_batch_size = 2000,
    kafka_flush_interval_ms = 500,
    kafka_num_consumers = 1,
    kafka_thread_per_consumer = 1;

CREATE TABLE IF NOT EXISTS analytics.conversions
(
    id String,
    action String,
    timestamp DateTime,
    click_id String,
    banner_id String,
    group_id String,
    campaign_id String,
    advertiser_id String,
    gaid String,
    oaid String,
    user_id String,
    stable_id String,
    bundle String,
    city String,
    country String,
    region String,
    price Decimal64(3),
    network String
)
ENGINE = MergeTree()
ORDER BY (timestamp, id)
TTL timestamp + INTERVAL 3 DAY DELETE;

CREATE MATERIALIZED VIEW IF NOT EXISTS analytics.conversions_parsed_mv
TO analytics.conversions AS
SELECT 
    JSONExtractString(raw_data, 'request_id') as id,
    parseDateTimeBestEffortOrZero(JSONExtractString(raw_data, 'timestamp')) as timestamp,
    JSONExtractString(raw_data, 'action') as action,
    JSONExtractString(raw_data, 'click_id') as click_id,
    JSONExtractString(raw_data, 'banner_id') as banner_id,
    JSONExtractString(raw_data, 'group_id') as group_id,
    JSONExtractString(raw_data, 'campaign_id') as campaign_id,
    JSONExtractString(raw_data, 'advertiser_id') as advertiser_id,
    JSONExtractString(raw_data, 'gaid') as gaid,
    JSONExtractString(raw_data, 'oaid') as oaid,
    JSONExtractString(raw_data, 'user_id') as user_id,
    JSONExtractString(raw_data, 'stable_id') as stable_id,
    JSONExtractString(raw_data, 'bundle') as bundle,
    JSONExtractString(raw_data, 'city') as city,
    JSONExtractString(raw_data, 'country') as country,
    JSONExtractString(raw_data, 'region') as region,
    if(empty(JSONExtractString(raw_data, 'price')), 0, 
       toDecimal64OrNull(JSONExtractString(raw_data, 'price'), 3)) as price,
    JSONExtractString(raw_data, 'network') as network
FROM analytics.kafka_conversions_raw;



------ win


CREATE TABLE IF NOT EXISTS analytics.kafka_wins_raw
(
    raw_data String
)
ENGINE = Kafka
SETTINGS
    kafka_broker_list = 'kafka1:9092,kafka2:9092,kafka3:9092',
    kafka_topic_list = 'win',
    kafka_group_name = 'clickhouse_win_consumer_v4',
    kafka_format = 'JSONAsString',
    kafka_max_block_size = 131072,
    kafka_poll_max_batch_size = 2000,
    kafka_flush_interval_ms = 500,
    kafka_num_consumers = 1;

CREATE TABLE IF NOT EXISTS analytics.wins
(
    id String,
    action String,
    timestamp DateTime,
    click_id String,
    banner_id String,
    group_id String,
    campaign_id String,
    advertiser_id String,
    gaid String,
    oaid String,
    user_id String,
    stable_id String,
    bundle String,
    city String,
    country String,
    region String,
    price Decimal64(3),
    network String
)
ENGINE = MergeTree()
ORDER BY (timestamp, id)
TTL timestamp + INTERVAL 3 DAY DELETE;

CREATE MATERIALIZED VIEW IF NOT EXISTS analytics.wins_parsed_mv
TO analytics.wins AS
SELECT 
    JSONExtractString(raw_data, 'request_id') as id,
    parseDateTimeBestEffortOrZero(JSONExtractString(raw_data, 'timestamp')) as timestamp,
    JSONExtractString(raw_data, 'action') as action,
    JSONExtractString(raw_data, 'click_id') as click_id,
    JSONExtractString(raw_data, 'banner_id') as banner_id,
    JSONExtractString(raw_data, 'group_id') as group_id,
    JSONExtractString(raw_data, 'campaign_id') as campaign_id,
    JSONExtractString(raw_data, 'advertiser_id') as advertiser_id,
    JSONExtractString(raw_data, 'gaid') as gaid,
    JSONExtractString(raw_data, 'oaid') as oaid,
    JSONExtractString(raw_data, 'user_id') as user_id,
    JSONExtractString(raw_data, 'stable_id') as stable_id,
    JSONExtractString(raw_data, 'bundle') as bundle,
    JSONExtractString(raw_data, 'city') as city,
    JSONExtractString(raw_data, 'country') as country,
    JSONExtractString(raw_data, 'region') as region,
    if(empty(JSONExtractString(raw_data, 'price')), 0, 
       toDecimal64OrNull(JSONExtractString(raw_data, 'price'), 3)) as price,
    JSONExtractString(raw_data, 'network') as network
FROM analytics.kafka_wins_raw;