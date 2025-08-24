-- +goose Up

CREATE TABLE IF NOT EXISTS `stamps` (
	`id` CHAR(36) NOT NULL,
	`name` VARCHAR(32) NOT NULL,
	`file_id` CHAR(36) NOT NULL,
	`creator_id` CHAR(36) NOT NULL,
	`is_unicode` BOOLEAN NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NOT NULL,
	`count_monthly`  INT UNSIGNED  NOT NULL DEFAULT 0,
	`count_total`  BIGINT UNSIGNED NOT NULL DEFAULT 0,
	PRIMARY KEY (`id`)
);
CREATE TABLE IF NOT EXISTS `tags` (
	`id` CHAR(36) NOT NULL,
	`name` VARCHAR(32) NOT NULL,
	`creator_id` CHAR(36) NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NOT NULL,
	PRIMARY KEY (`id`),
	UNIQUE KEY (`name`)
);
CREATE TABLE IF NOT EXISTS `stamp_daily_usages` (
	`stamp_id` CHAR(36) NOT NULL,
	`date` DATE NOT NULL,
	`reaction_count` INT UNSIGNED NOT NULL,
	`message_count` INT UNSIGNED NOT NULL,
	PRIMARY KEY (`stamp_id`, `date`),
	FOREIGN KEY (`stamp_id`) REFERENCES `stamps`(`id`)
);
CREATE TABLE IF NOT EXISTS `stamp_descriptions` (
	`stamp_id` CHAR(36) NOT NULL,
	`description` TEXT NOT NULL,
	`creator_id` CHAR(36) NOT NULL,
	`created_at` DATETIME NOT NULL,
	`updated_at` DATETIME NOT NULL,
	PRIMARY KEY (`stamp_id`,`creator_id`),
	FOREIGN KEY (`stamp_id`) REFERENCES `stamps`(`id`)
);
CREATE TABLE IF NOT EXISTS `stamp_tags` (
	`stamp_id` CHAR(36) NOT NULL,
	`tag_id` CHAR(36) NOT NULL,
	`creator_id` CHAR(36) NOT NULL,
	PRIMARY KEY (`stamp_id`, `tag_id`),
	FOREIGN KEY (`stamp_id`) REFERENCES `stamps`(`id`),
	FOREIGN KEY (`tag_id`) REFERENCES `tags`(`id`) ON DELETE CASCADE
);