DROP TABLE IF EXISTS transcripts;

CREATE TABLE IF NOT EXISTS transcripts
(
    id            integer primary key,
    text          text,
    role          varchar(255),
    is_valid      boolean default true,
    simulation_id int,
    FOREIGN KEY (simulation_id) REFERENCES simulations (id)
);
