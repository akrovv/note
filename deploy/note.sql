DROP TABLE IF EXISTS `note`;
CREATE TABLE `note` (
    `id` INT(11) PRIMARY KEY AUTO_INCREMENT,
    `text` TEXT,
    `created_at` DATETIME,
    `updated_at` DATETIME
) ENGINE=InnoDB DEFAULT CHARSET=utf8;