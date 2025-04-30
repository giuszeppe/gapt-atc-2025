package stores

import (
	"database/sql"
)

type User struct {
	ID       int
	Username string
	Password string
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) Store(user User) error {
	// _, err := s.db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	return nil
}

func (s *UserStore) View() ([]User, error) {
	// rows, err := s.db.Query("SELECT id, name, email FROM users")
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()
	//
	// var users []User
	// for rows.Next() {
	// 	var user User
	// 	if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
	// 		return nil, err
	// 	}
	// 	users = append(users, user)
	// }
	//
	// return users, nil
	return []User{}, nil
}

func (s *UserStore) Exist(user User) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", user.Username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *UserStore) GetUserWithUsername(username string) (User, error) {
	var password string
	err := s.db.QueryRow("SELECT password FROM users WHERE username=?", username).Scan(&password)

	if err != nil {
        return User{}, err
	}

	return User{
		Username: username,
		Password: password,
	}, nil

}
