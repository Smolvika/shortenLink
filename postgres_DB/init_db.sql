CREATE TABLE links
(
    short_url   VARCHAR(10)   PRIMARY KEY,
    original_url VARCHAR(1024) NOT NULL,
    expiration_date VARCHAR(11) NOT NULL
);