CREATE DATABASE db1;
USE db1;
CREATE TABLE requests (
    id int(2) NOT NULL AUTO_INCREMENT,
    name varchar(127) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
    email varchar(95) COLLATE utf8_unicode_ci NOT NULL,
    url longtext COLLATE utf8_unicode_ci NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE articles (
      id int(4) NOT NULL AUTO_INCREMENT,
      hid varchar(127) NOT NULL,
      request_id int(2) NOT NULL,
      publisher varchar(127) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
      location varchar(255) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
      title varchar(255) COLLATE utf8_unicode_ci NOT NULL,
      href longtext COLLATE utf8_unicode_ci NOT NULL,
      image longtext COLLATE utf8_unicode_ci NOT NULL,
      features longtext COLLATE utf8_unicode_ci NOT NULL,
      price int NOT NULL,
      last_update varchar(127) COLLATE utf8_unicode_ci NOT NULL DEFAULT '',
      PRIMARY KEY (id),
      CONSTRAINT FK_REQUEST_ARTICLE FOREIGN KEY (request_id) REFERENCES requests (id)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;