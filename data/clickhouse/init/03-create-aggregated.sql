CREATE TABLE IF NOT EXISTS analytics.requests_hour
(
    action String,
    timestamp DateTime,
    banner_id String,
    group_id String,
    campaign_id String,
    advertiser_id String,
    city String,
    country String,
    region String,
    price Decimal64(3),
    count Integer,
    network String,
    bundle String
)
ENGINE = MergeTree()
ORDER BY (timestamp);

CREATE TABLE IF NOT EXISTS analytics.impressions_hour
(
    action String,
    timestamp DateTime,
    banner_id String,
    group_id String,
    campaign_id String,
    advertiser_id String,
    city String,
    country String,
    region String,
    price Decimal64(3),
    count Integer,
    network String,
    bundle String
)
ENGINE = MergeTree()
ORDER BY (timestamp);

CREATE TABLE IF NOT EXISTS analytics.responses_hour
(
    action String,
    timestamp DateTime,
    banner_id String,
    group_id String,
    campaign_id String,
    advertiser_id String,
    city String,
    country String,
    region String,
    price Decimal64(3),
    count Integer,
    network String,
    bundle String
)
ENGINE = MergeTree()
ORDER BY (timestamp);

CREATE TABLE IF NOT EXISTS analytics.clicks_hour
(
    action String,
    timestamp DateTime,
    banner_id String,
    group_id String,
    campaign_id String,
    advertiser_id String,
    city String,
    country String,
    region String,
    price Decimal64(3),
    count Integer,
    network String,
    bundle String
)
ENGINE = MergeTree()
ORDER BY (timestamp);

CREATE TABLE IF NOT EXISTS analytics.conversions_hour
(
    action String,
    timestamp DateTime,
    banner_id String,
    group_id String,
    campaign_id String,
    advertiser_id String,
    city String,
    country String,
    region String,
    price Decimal64(3),
    count Integer,
    network String,
    bundle String
)
ENGINE = MergeTree()
ORDER BY (timestamp);

CREATE TABLE IF NOT EXISTS analytics.wins_hour
(
    action String,
    timestamp DateTime,
    banner_id String,
    group_id String,
    campaign_id String,
    advertiser_id String,
    city String,
    country String,
    region String,
    price Decimal64(3),
    count Integer,
    network String,
    bundle String
)
ENGINE = MergeTree()
ORDER BY (timestamp);