-- +goose Up
-- 新しい `count` カラムを追加
ALTER TABLE `stamps` ADD COLUMN `count` BIGINT UNSIGNED NOT NULL DEFAULT 0;
ALTER TABLE `stamps` DROP COLUMN `count_monthly`;

-- +goose Down
-- ロールバック時に元の状態に戻す
ALTER TABLE `stamps` ADD COLUMN `count_monthly` BIGINT UNSIGNED NOT NULL DEFAULT 0;
ALTER TABLE `stamps` DROP COLUMN `count`;