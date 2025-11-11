-- +goose Up
-- 既存の `count_monthly` カラムを削除
ALTER TABLE `stamps` DROP COLUMN `count_monthly`;

-- 新しい `count` カラムを追加
ALTER TABLE `stamps` ADD COLUMN `count` BIGINT UNSIGNED NOT NULL DEFAULT 0;
-- +goose Down