CREATE DATABASE IF NOT EXISTS key_management
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE key_management;

CREATE TABLE IF NOT EXISTS configs (
  `key` VARCHAR(100) NOT NULL,
  `value` TEXT NOT NULL,
  `description` VARCHAR(255) DEFAULT '',
  `created_at` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` DATETIME(3) NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO configs (`key`, `value`, `description`)
VALUES ('copy_template', '{{key}}', 'Global copy template for key values')
ON DUPLICATE KEY UPDATE
  `description` = VALUES(`description`);
