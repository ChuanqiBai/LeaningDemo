package service

import "sync"

type UserStore interface {
	Save(user *User) error
	Find(userName string) (*User, error)
}

type InMemoryUserStore struct {
	mutex sync.RWMutex
	users map[string]*User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*User),
	}
}

func (store *InMemoryUserStore) Save(user *User) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.users[user.UserName] != nil {
		return ErrAlreadyExists
	}

	store.users[user.UserName] = user.Clone()
	return nil
}

func (store *InMemoryUserStore) Find(userName string) (*User, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	user := store.users[userName]
	if user == nil {
		return nil, nil
	}
	return user.Clone(), nil
}
