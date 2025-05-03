DROP TABLE IF EXISTS scenarios;

CREATE TABLE IF NOT EXISTS scenarios  (
    id integer primary key,
    type varchar(255),
    name varchar(255)
);
