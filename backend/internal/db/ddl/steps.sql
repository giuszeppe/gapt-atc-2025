
DROP TABLE IF EXISTS steps;
DROP TABLE IF EXISTS extended_steps;

CREATE TABLE IF NOT EXISTS steps (
     id integer primary key,
    idx integer default 1,
    text TEXT,
    role varchar(200),
    scenario_id integer not null,
    foreign key(scenario_id) references scenarios(id)
);

CREATE TABLE IF NOT EXISTS extended_steps (
    id integer primary key,
    idx integer default 1,
    text TEXT,
    role varchar(200),
    scenario_id integer not null,
    foreign key(scenario_id) references scenarios(id)
);
