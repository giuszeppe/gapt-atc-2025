DROP TABLE IF EXISTS scenarios;
DROP INDEX IF EXISTS idx_scenarios_type;

CREATE TABLE IF NOT EXISTS scenarios
(
    id   integer primary key,
    type varchar(255),
    name varchar(255) UNIQUE
);

CREATE INDEX IF NOT EXISTS idx_scenarios_type ON scenarios (type);
