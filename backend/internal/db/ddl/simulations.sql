DROP TABLE IF EXISTS simulations;
DROP INDEX IF EXISTS idx_simulations_scenario_id;
DROP INDEX IF EXISTS idx_simulations_lobby_id;

CREATE TABLE IF NOT EXISTS simulations
(
    id                          integer PRIMARY KEY,  -- Auto-incrementing ID
    scenario_id                 INT         NOT NULL,
    input_type                  VARCHAR(50) NOT NULL, -- block, text, speech
    scenario_type               VARCHAR(50) NOT NULL, -- takeoff, en route, landing
    simulation_advancement_type VARCHAR(50) NOT NULL,-- continuous, steps
    mode                        VARCHAR(50) NOT NULL, -- single, multi
    tower_user_id               INT,
    aircraft_user_id            INT,
    lobby_id                    varchar(6) UNIQUE,
    FOREIGN KEY (tower_user_id) REFERENCES users (id),
    FOREIGN KEY (aircraft_user_id) REFERENCES users (id),
    FOREIGN KEY (scenario_id) REFERENCES scenarios (id)
);

CREATE INDEX IF NOT EXISTS idx_simulations_scenario_id ON simulations (scenario_id);
CREATE INDEX IF NOT EXISTS idx_simulations_lobby_id ON simulations (lobby_id);


