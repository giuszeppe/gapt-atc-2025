package stores

import (
	"database/sql"
	"fmt"
)

type Scenario struct {
	ID            int
	Name          string
	Type          string
	Steps         []Step
	ExtendedSteps []Step
}

type Step struct {
	Index int
	Text  string
	Role  string
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
		var step Step

		if err := rows.Scan(&scenario.ID, &scenario.Name, &scenario.Type); err != nil {
			return []Scenario{}, err
		}
		// Check if scenario is already in map
		fmt.Println(scenario, step)
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

    res := [][]Step{[]Step{}, []Step{}}

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
