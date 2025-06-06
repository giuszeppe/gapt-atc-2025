package stores

import (
	"database/sql"
	"errors"
)

type Transcript struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Id      int    `json:"id"`
	Text    string `json:"content"`
	Role    string `json:"role"` // 'tower', 'aircraft'
	IsValid bool   `json:"is_valid"`
}

type Simulation struct {
	// Role                      string `json:"role"`                        // tower, aircraft
	Id                        int        `json:"id"`
	ScenarioId                int        `json:"scenario_id"`
	InputType                 string     `json:"input_type"`                  // block, text, speech
	ScenarioType              string     `json:"scenario_type"`               // takeoff, enroute, landing
	SimulationAdvancementType string     `json:"simulation_advancement_type"` // continuous, steps
	Mode                      string     `json:"mode"`                        // single, multi
	Transcript                Transcript `json:"transcript,omitempty"`
	TowerUserId               int        `json:"tower_user_id"`
	AircraftUserId            int        `json:"aircraft_user_id"`
	LobbyCode                 string     `json:"lobby_code,omitempty"`
}

type Scenario struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Steps         []Step `json:"steps,omitempty"`
	ExtendedSteps []Step `json:"extended_steps,omitempty"`
}

type Step struct {
	Index int    `json:"index"`
	Text  string `json:"text"`
	Role  string `json:"role"`
}

type ScenarioStore struct {
	db *sql.DB
}

func NewScenarioStore(db *sql.DB) *ScenarioStore {
	return &ScenarioStore{db: db}
}

func (s *ScenarioStore) View(scenarioType string) ([]Scenario, error) {

	stmt, err := s.db.Prepare(`SELECT sc.id, sc.name, sc.type
        FROM scenarios sc
        WHERE sc.type = ?;
        `)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(scenarioType)
	if err != nil {
		return []Scenario{}, err
	}

	scenarioMap := make(map[int]*Scenario)
	for rows.Next() {
		var scenario Scenario

		if err := rows.Scan(&scenario.ID, &scenario.Name, &scenario.Type); err != nil {
			return []Scenario{}, err
		}
		// Check if scenario is already in map
		if _, exists := scenarioMap[scenario.ID]; !exists {
			scenarioMap[scenario.ID] = &scenario
		}
	}
	// Convert map to slice
	var scenarios []Scenario
	for _, scenario := range scenarioMap {
		scenarios = append(scenarios, *scenario)
	}

	return scenarios, nil
}

func (s *ScenarioStore) GetScenarioStepsForId(scenarioId int) ([][]Step, error) {
	stmt, err := s.db.Prepare(`
        SELECT st.idx, st.text, st.role, est.idx, est.text, est.role
        FROM scenarios s
        LEFT JOIN steps st on s.id=st.scenario_id
        LEFT JOIN extended_steps est on s.id=est.scenario_id AND est.idx = st.idx
        WHERE s.id= ?;
        `)
	if err != nil {
		return [][]Step{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(scenarioId)
	if err != nil {
		return [][]Step{}, err
	}
	defer rows.Close()

	res := [][]Step{{}, {}}

	for rows.Next() {
		var step Step
		var extendedStep Step
		if err := rows.Scan(&step.Index, &step.Text, &step.Role, &extendedStep.Index, &extendedStep.Text, &extendedStep.Role); err != nil {
			return [][]Step{}, err
		}
		res[0] = append(res[0], step)
		res[1] = append(res[1], extendedStep)

	}

	return res, nil
}

func (s *ScenarioStore) StoreSimulation(scenarioId, userId int, role, inputType, scenarioType, advancementType, mode, lobbyCode string) (Simulation, error) {
	// Example: userId is assumed to be both tower and aircraft user for simplicity
	query := `
		INSERT INTO simulations (
			scenario_id,
			input_type,
			scenario_type,
			simulation_advancement_type,
			mode,
			tower_user_id,
			aircraft_user_id,
			lobby_id	
		) VALUES ($1, $2, $3, $4, $5, $6, $7,$8)
		RETURNING id;
	`
	towerId := -1
	aircraftId := -1
	if role == "tower" {
		towerId = userId
	} else {
		aircraftId = userId
	}

	var id int
	var lobbyValue *string
	lobbyValue = &lobbyCode
	if lobbyCode == "" {
		lobbyValue = nil
	}
	err := s.db.QueryRow(
		query,
		scenarioId,
		inputType,
		scenarioType,
		advancementType,
		mode,
		towerId, // assuming same user for tower and aircraft
		aircraftId,
		lobbyValue,
	).Scan(&id)

	if err != nil {
		return Simulation{}, err
	}

	simulation := Simulation{
		Id:                        id,
		ScenarioId:                scenarioId,
		InputType:                 inputType,
		ScenarioType:              scenarioType,
		SimulationAdvancementType: advancementType,
		Mode:                      mode,
		TowerUserId:               towerId,
		AircraftUserId:            aircraftId,
	}
	if lobbyCode != "" {
		simulation.LobbyCode = lobbyCode
	}

	return simulation, nil
}

func (s *ScenarioStore) AddTranscriptToSimulationUsingLobbyCode(lobbyCode string, messages []Message) error {
	query := `SELECT id FROM simulations WHERE lobby_id = $1;`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var simulationId int
	err = stmt.QueryRow(lobbyCode).Scan(&simulationId)
	if err != nil {
		return err
	}

	err = s.addMessagesToSimulation(simulationId, messages)
	if err != nil {
		return err
	}

	return nil
}

func (s *ScenarioStore) DoesLobbyCodeExist(code string) (bool, error) {
	query := "SELECT count(*) FROM simulations WHERE lobby_id = ?;"
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	var count int
	err = stmt.QueryRow(code).Scan(&count)
	if err != nil {
		return false, err
	}

	return count != 0, nil
}

func (s *ScenarioStore) addMessagesToSimulation(simulationId int, messages []Message) error {
	query := `INSERT INTO transcripts (text,role,simulation_id, is_valid) VALUES`
	values := []any{}

	for idx, message := range messages {
		values = append(values, message.Text, message.Role, simulationId, message.IsValid)
		if idx == 0 {
			query += `(?,?,?,?)`
		} else {
			query += `,(?,?,?,?)`
		}
	}
	query += ";"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}

	return nil
}

// EndSimulation adds transcript to simulation
func (s *ScenarioStore) EndSimulation(scenarioId int, messages []Message) error {
	return s.addMessagesToSimulation(scenarioId, messages)
}

/**
 -type
	-scenario_specific
		-transcript 1
		- transcript 2
*/

func (s *ScenarioStore) GetGroupedTranscripts() (map[string]map[string]map[int]*Transcript, error) {
	query := `SELECT t.id, t.text, t.role, s.name, s.type, t.simulation_id, t.is_valid FROM transcripts t
    LEFT JOIN simulations ON simulations.id = t.simulation_id
    LEFT JOIN main.scenarios s on simulations.scenario_id = s.id`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transcripts := make(map[string]map[string]map[int]*Transcript)
	for rows.Next() {
		var message Message
		var scenarioType string
		var scenarioName string // here I am assuming scenario names are unique
		var simulationId int
		if err := rows.Scan(&message.Id, &message.Text, &message.Role, &scenarioName, &scenarioType, &simulationId, &message.IsValid); err != nil {
			return nil, err
		}
		if _, ok := transcripts[scenarioType]; !ok {
			transcripts[scenarioType] = make(map[string]map[int]*Transcript)
		}
		if _, ok := transcripts[scenarioType][scenarioName]; !ok {
			transcripts[scenarioType][scenarioName] = make(map[int]*Transcript)
		}
		if _, ok := transcripts[scenarioType][scenarioName][simulationId]; !ok {
			transcripts[scenarioType][scenarioName][simulationId] = &Transcript{}
		}
		transcripts[scenarioType][scenarioName][simulationId].Messages = append(transcripts[scenarioType][scenarioName][simulationId].Messages, message)
	}
	return transcripts, nil
}
func (s *ScenarioStore) GetTranscriptBySimulationId(simulationId int) (Transcript, error) {
	query := `SELECT id,text,role, is_valid FROM transcripts WHERE simulation_id = ?`

	messages := []Message{}
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return Transcript{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(simulationId)
	if err != nil {
		return Transcript{}, err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		found = true
		var message Message
		if err := rows.Scan(&message.Id, &message.Text, &message.Role, &message.IsValid); err != nil {
			return Transcript{}, err
		}
		messages = append(messages, message)
	}
	if !found {
		return Transcript{}, errors.New("Not transcripts found")
	}
	return Transcript{Messages: messages}, nil

}

func (s *ScenarioStore) GetUserRoleForSimulation(simulationId int) (string, error) {
	query := `SELECT simulations.tower_user_id, simulations.aircraft_user_id FROM simulations WHERE id = ?`
	rows, err := s.db.Query(query, simulationId)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		var towerUserId, aircraftUserId int
		if err := rows.Scan(&towerUserId, &aircraftUserId); err != nil {
			return "", err
		}
		if towerUserId != -1 {
			return "aircraft", nil
		} else {
			return "tower", nil
		}
	}
	return "", nil
}

func (s *ScenarioStore) GetSimulationByLobbyCode(lobbyCode string) (Simulation, error) {
	query := `SELECT 
		s.id, 
		s.scenario_id, 
		s.input_type, 
		s.scenario_type, 
		s.simulation_advancement_type, 
		s.mode, 
		s.tower_user_id, 
		s.aircraft_user_id 
	FROM simulations s 
	WHERE lobby_id = ?`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return Simulation{}, err
	}
	defer stmt.Close()
	var simulation Simulation
	err = stmt.QueryRow(lobbyCode).Scan(
		&simulation.Id,
		&simulation.ScenarioId,
		&simulation.InputType,
		&simulation.ScenarioType,
		&simulation.SimulationAdvancementType,
		&simulation.Mode,
		&simulation.TowerUserId,
		&simulation.AircraftUserId,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Simulation{}, errors.New("simulation not found")
		}
		return Simulation{}, err
	}
	return simulation, nil
}

func (s *ScenarioStore) UpdateSimulationRoleIds(simulationId, userId int, role string) error {
	var query string
	if role == "tower" {
		query = `UPDATE simulations SET tower_user_id = ? WHERE id = ?`
	} else if role == "aircraft" {
		query = `UPDATE simulations SET aircraft_user_id = ? WHERE id = ?`
	}
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(userId, simulationId)
	if err != nil {
		return err
	}
	return nil
}

func (s *ScenarioStore) GetScenarioStepsForSimulationId(simulationId int) ([]Step, error) {
	query := `
		SELECT est.idx, est.text, est.role
		FROM simulations s
		LEFT JOIN extended_steps est on s.scenario_id=est.scenario_id 
		WHERE s.id= ?;`
	stmt, err := s.db.Prepare(query)
	if err != nil {
		return []Step{}, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(simulationId)
	if err != nil {
		return []Step{}, err
	}
	defer rows.Close()

	res := []Step{}

	for rows.Next() {
		var extendedStep Step
		if err := rows.Scan(&extendedStep.Index, &extendedStep.Text, &extendedStep.Role); err != nil {
			return []Step{}, err
		}
		res = append(res, extendedStep)

	}

	return res, nil
}
