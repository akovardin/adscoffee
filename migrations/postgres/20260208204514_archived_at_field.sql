-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

-- Добавляем поле archived_at в таблицу advertisers
ALTER TABLE public.advertisers
ADD COLUMN archived_at timestamp with time zone;

-- Добавляем поле archived_at в таблицу campaigns
ALTER TABLE public.campaigns
ADD COLUMN archived_at timestamp with time zone;

-- Добавляем поле archived_at в таблицу bgroups
ALTER TABLE public.bgroups
ADD COLUMN archived_at timestamp with time zone;

-- Добавляем поле archived_at в таблицу banners
ALTER TABLE public.banners
ADD COLUMN archived_at timestamp with time zone;

-- Добавляем поле archived_at в таблицу networks
ALTER TABLE public.networks
ADD COLUMN archived_at timestamp with time zone;

-- Создаем индексы для оптимизации запросов с archived_at
CREATE INDEX idx_advertisers_archived_at ON public.advertisers USING btree (archived_at);
CREATE INDEX idx_campaigns_archived_at ON public.campaigns USING btree (archived_at);
CREATE INDEX idx_bgroups_archived_at ON public.bgroups USING btree (archived_at);
CREATE INDEX idx_banners_archived_at ON public.banners USING btree (archived_at);
CREATE INDEX idx_networks_archived_at ON public.networks USING btree (archived_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

-- Удаляем индексы для archived_at
DROP INDEX IF EXISTS idx_advertisers_archived_at;
DROP INDEX IF EXISTS idx_campaigns_archived_at;
DROP INDEX IF EXISTS idx_bgroups_archived_at;
DROP INDEX IF EXISTS idx_banners_archived_at;
DROP INDEX IF EXISTS idx_networks_archived_at;


ALTER TABLE public.networks
DROP COLUMN IF EXISTS archived_at;

ALTER TABLE public.banners
DROP COLUMN IF EXISTS archived_at;

ALTER TABLE public.bgroups
DROP COLUMN IF EXISTS archived_at;

ALTER TABLE public.campaigns
DROP COLUMN IF EXISTS archived_at;

ALTER TABLE public.advertisers
DROP COLUMN IF EXISTS archived_at;
-- +goose StatementEnd