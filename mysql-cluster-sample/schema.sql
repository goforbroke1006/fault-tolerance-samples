CREATE TABLE IF NOT EXISTS `employee` (
  `id`      int(6) unsigned NOT NULL AUTO_INCREMENT,
  `name`    varchar(128)    NOT NULL,
  `city`    varchar(128)    NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4;