DROP TABLE IF EXISTS simulations;
CREATE TABLE IF NOT EXISTS simulations (
    id integer PRIMARY KEY, -- Auto-incrementing ID scenario_id INT NOT NULL,
    input_type VARCHAR(50) NOT NULL,                 -- block, text, speech
    scenario_type VARCHAR(50) NOT NULL,              -- takeoff, en route, landing
    simulation_advancement_type VARCHAR(50) NOT NULL,-- continuous, steps
    mode VARCHAR(50) NOT NULL,                       -- single, multi
    tower_user_id INT,
    aircraft_user_id INT,
    scenario_id int,
    lobby_id varchar(6)
);

