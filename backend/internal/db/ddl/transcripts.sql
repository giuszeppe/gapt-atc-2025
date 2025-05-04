DROP TABLE IF EXISTS transcripts;

CREATE TABLE IF NOT EXISTS transcripts (
    id int primary key,
    text text,
    role varchar(255),
    simulation_id int
);
