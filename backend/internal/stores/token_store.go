package stores

import "errors"

type TokenStore struct {
	items map[string]User
	zz    chan User
}

func NewTokenStore() *TokenStore {
	m := TokenStore{
		items: make(map[string]User),
		zz:    make(chan User),
	}
	go m.start()
	return &m
}

func (m *TokenStore) start() {
	for {
		select {
		case user := <-m.zz:
			m.items[user.Token] = user
		}
	}
}

func (m *TokenStore) Store(x User) error {
	m.zz <- x
	return nil
}

func (m *TokenStore) View() (map[string]User, error) {
	ret := make(map[string]User)
	for k, v := range m.items {
		ret[k] = v
	}
	return ret, nil
}

func (m *TokenStore) Exist(x User) (bool, error) {
	_, ok := m.items[x.Token]
	return ok, nil
}

func (m *TokenStore) GetUserByToken(t string) (User, error) {
	usr, ok := m.items[t]
	if !ok {
		return usr, errors.New("user not foundddd")
	}
	return usr, nil
}
