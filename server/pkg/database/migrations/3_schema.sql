-- +goose Up
-- tagsテーブルのnameカラムにUNIQUE制約を追加
ALTER TABLE `tags` ADD UNIQUE (`name`);
