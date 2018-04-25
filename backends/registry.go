package backends

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrClientAlreadyExists = errors.New("client already exists")
)

var (
	clients      = make(map[string]*client)
	clientsMutex sync.RWMutex
)

type NewClientFunc func(c *Config) (Client, error)

type client struct {
	name          string
	priority      int
	newClientFunc NewClientFunc
}

func Register(name string, p int, f NewClientFunc) error {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()
	if _, ok := clients[name]; ok {
		return ErrClientAlreadyExists
	}
	clients[name] = &client{
		name:          name,
		priority:      p,
		newClientFunc: f,
	}
	return nil
}

func clientByName(name string) (*client, error) {
	clientsMutex.RLock()
	defer clientsMutex.RUnlock()
	c, ok := clients[name]
	if !ok {
		return nil, fmt.Errorf("no such backend client '%s'", name)
	}
	return c, nil
}

func defaultClient() *client {
	clientsMutex.RLock()
	defer clientsMutex.RUnlock()
	if len(clients) == 0 {
		// This will never happen
		panic("no clients registered")
	}
	var best *client
	for _, c := range clients {
		if best == nil {
			best = c
			continue
		}
		if c.priority > best.priority {
			best = c
		}
	}
	return best
}
