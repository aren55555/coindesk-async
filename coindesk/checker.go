package coindesk

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	pollingInterval = 1 * time.Second
	apiURL          = "https://api.coindesk.com/v1/bpi/currentprice.json"
)

// Checker presents an interface for fetching values from the CoinDesk api on a repeated interval.
type Checker struct {
	started bool
	stop    chan struct{}
	mu      sync.RWMutex
	current *Value
}

// New will prepare a new instance of a Checker.
func New() *Checker {
	return &Checker{
		stop: make(chan struct{}),
	}
}

// GetValue will return the Checker's current Value if there was one stored.
func (c *Checker) GetValue() (Value, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.current == nil {
		return Value{}, errors.New("no value set")
	}
	return *c.current, nil
}

// Stop will cease the polling.
// Stop will cease the polling.
func (c *Checker) Stop() {
	fmt.Println("stopping...")
	c.stop <- struct{}{}
}

// Start will start/resume the polling.
func (c *Checker) Start() {
	fmt.Println("starting...")
	go c.do()
}

func (c *Checker) do() {
	ticker := time.NewTicker(pollingInterval)
	for {
		select {
		case <-ticker.C:
			go c.fetchAndUpdate()
		case <-c.stop:
			ticker.Stop()
			return
		}
	}
}

func (c *Checker) fetchAndUpdate() {
	fmt.Println("fetching a new price...")

	request, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return
	}
	fetchTime := time.Now()
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return
	}
	defer response.Body.Close()

	responseData := Response{}
	if err := json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return
	}

	newValue := &Value{
		FetchTime: fetchTime,
		Result:    responseData,
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.current = newValue
}
