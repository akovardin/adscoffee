-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE public.advertisers (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    title text,
    info text,
    active boolean,
    start timestamp with time zone,
    "end" timestamp with time zone,
    targeting text,
    budget text,
    capping text,
    timetable text,
    ord_contract text,
    ord_enable boolean
);

CREATE TABLE public.campaigns (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    title text,
    active boolean,
    bundle text,
    start timestamp with time zone,
    "end" timestamp with time zone,
    targeting text,
    budget text,
    capping text,
    timetable text,
    advertiser_id bigint
);

CREATE SEQUENCE public.campaigns_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.campaigns_id_seq OWNED BY public.campaigns.id;

CREATE SEQUENCE public.advertisers_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

CREATE TABLE public.bgroups (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    title text,
    active boolean,
    price bigint,
    start timestamp with time zone,
    "end" timestamp with time zone,
    targeting text,
    budget text,
    capping text,
    timetable text,
    campaign_id bigint
);

CREATE SEQUENCE public.bgroups_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.bgroups_id_seq OWNED BY public.bgroups.id;

CREATE TABLE public.banners (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    title text,
    label text,
    description text,
    active boolean,
    price bigint,
    image text,
    icon text,
    start timestamp with time zone,
    "end" timestamp with time zone,
    clicktracker text,
    imptracker text,
    target text,
    targeting text,
    budget text,
    capping text,
    bgroup_id bigint,
    timetable text,
    erid text,
    contract text,
    ord_category text,
    ord_targeting text,
    ord_format text,
    ord_kktu text,
    expected_win_rate numeric
);

CREATE SEQUENCE public.banners_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.banners_id_seq OWNED BY public.banners.id;


CREATE TABLE public.networks (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    title text,
    name text
);

CREATE SEQUENCE public.networks_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public.networks_id_seq OWNED BY public.networks.id;

CREATE TABLE public.users (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    name text,
    account text,
    password character varying(60),
    pass_updated_at text,
    login_retry_count bigint,
    locked boolean,
    locked_at timestamp with time zone,
    reset_password_token text,
    reset_password_token_created_at timestamp with time zone,
    reset_password_token_expired_at timestamp with time zone,
    totp_secret text,
    is_totp_setup boolean,
    last_used_totp_code text,
    last_totp_code_used_at timestamp with time zone,
    session_secure character varying(32)
);

CREATE SEQUENCE public.users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;



CREATE TABLE public.media_libraries (
    id bigint NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    selected_type text,
    file text,
    user_id bigint,
    folder boolean DEFAULT false,
    parent_id bigint DEFAULT 0
);

CREATE SEQUENCE public.media_libraries_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.media_libraries_id_seq OWNED BY public.media_libraries.id;

ALTER TABLE ONLY public.advertisers ALTER COLUMN id SET DEFAULT nextval('public.advertisers_id_seq'::regclass);

ALTER TABLE ONLY public.banners ALTER COLUMN id SET DEFAULT nextval('public.banners_id_seq'::regclass);

ALTER TABLE ONLY public.bgroups ALTER COLUMN id SET DEFAULT nextval('public.bgroups_id_seq'::regclass);

ALTER TABLE ONLY public.campaigns ALTER COLUMN id SET DEFAULT nextval('public.campaigns_id_seq'::regclass);

ALTER TABLE ONLY public.media_libraries ALTER COLUMN id SET DEFAULT nextval('public.media_libraries_id_seq'::regclass);

ALTER TABLE ONLY public.networks ALTER COLUMN id SET DEFAULT nextval('public.networks_id_seq'::regclass);

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);

ALTER TABLE ONLY public.advertisers
    ADD CONSTRAINT advertisers_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.banners
    ADD CONSTRAINT banners_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.bgroups
    ADD CONSTRAINT bgroups_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.campaigns
    ADD CONSTRAINT campaigns_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.media_libraries
    ADD CONSTRAINT media_libraries_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.networks
    ADD CONSTRAINT networks_pkey PRIMARY KEY (id);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);

CREATE INDEX idx_advertisers_deleted_at ON public.advertisers USING btree (deleted_at);

CREATE INDEX idx_banners_deleted_at ON public.banners USING btree (deleted_at);

CREATE INDEX idx_bgroups_deleted_at ON public.bgroups USING btree (deleted_at);

CREATE INDEX idx_campaigns_deleted_at ON public.campaigns USING btree (deleted_at);

CREATE INDEX idx_media_libraries_deleted_at ON public.media_libraries USING btree (deleted_at);

CREATE INDEX idx_media_libraries_parent_id ON public.media_libraries USING btree (parent_id);

CREATE INDEX idx_media_libraries_user_id ON public.media_libraries USING btree (user_id);

CREATE INDEX idx_networks_deleted_at ON public.networks USING btree (deleted_at);

CREATE UNIQUE INDEX idx_users_account ON public.users USING btree (account) WHERE ((account <> ''::text) AND (deleted_at IS NULL));

CREATE INDEX idx_users_deleted_at ON public.users USING btree (deleted_at);

CREATE UNIQUE INDEX idx_users_reset_password_token ON public.users USING btree (reset_password_token) WHERE (reset_password_token <> ''::text);

ALTER TABLE ONLY public.banners
    ADD CONSTRAINT fk_banners_bgroup FOREIGN KEY (bgroup_id) REFERENCES public.bgroups(id);

ALTER TABLE ONLY public.bgroups
    ADD CONSTRAINT fk_bgroups_campaign FOREIGN KEY (campaign_id) REFERENCES public.campaigns(id);

ALTER TABLE ONLY public.campaigns
    ADD CONSTRAINT fk_campaigns_advertiser FOREIGN KEY (advertiser_id) REFERENCES public.advertisers(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP INDEX IF EXISTS idx_advertisers_deleted_at;

DROP INDEX IF EXISTS idx_banners_deleted_at;

DROP INDEX IF EXISTS idx_bgroups_deleted_at;

DROP INDEX IF EXISTS idx_campaigns_deleted_at;

DROP INDEX IF EXISTS idx_media_libraries_deleted_at;

DROP INDEX IF EXISTS idx_media_libraries_parent_id;

DROP INDEX IF EXISTS idx_media_libraries_user_id;

DROP INDEX IF EXISTS idx_networks_deleted_at;

DROP INDEX IF EXISTS idx_users_account;

DROP INDEX IF EXISTS idx_users_deleted_at;

DROP INDEX IF EXISTS idx_users_reset_password_token;

DROP TABLE IF EXISTS public.media_libraries;

DROP SEQUENCE IF EXISTS public.media_libraries_id_seq;

DROP TABLE IF EXISTS public.users;

DROP SEQUENCE IF EXISTS public.users_id_seq;

DROP TABLE IF EXISTS public.networks;

DROP SEQUENCE IF EXISTS public.networks_id_seq;

DROP TABLE IF EXISTS public.banners;

DROP SEQUENCE IF EXISTS public.banners_id_seq;

DROP TABLE IF EXISTS public.bgroups;

DROP SEQUENCE IF EXISTS public.bgroups_id_seq;

DROP TABLE IF EXISTS public.campaigns;

DROP SEQUENCE IF EXISTS public.campaigns_id_seq;

DROP TABLE IF EXISTS public.advertisers;

DROP SEQUENCE IF EXISTS public.advertisers_id_seq;
-- +goose StatementEnd
