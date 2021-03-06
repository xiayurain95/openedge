package api_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/baidu/openedge/api"
	"github.com/baidu/openedge/config"
	"github.com/baidu/openedge/module"
	"github.com/baidu/openedge/trans/http"
	"github.com/stretchr/testify/assert"
)

type mockEngine struct {
	pass bool
}

func (e *mockEngine) Authenticate(username, password string) bool {
	fmt.Println("Authenticate")
	return e.pass
}

func (e *mockEngine) Start(_ config.Module) error {
	fmt.Println("Start")
	return nil
}

func (e *mockEngine) Restart(_ string) error {
	fmt.Println("restart")
	return nil
}

func (e *mockEngine) Stop(_ string) error {
	fmt.Println("Stop")
	return nil
}

func TestAPIHttp(t *testing.T) {
	s, err := api.NewServer(&mockEngine{pass: true}, http.ServerConfig{Address: "tcp://127.0.0.1:0", Timeout: time.Minute})
	assert.NoError(t, err)
	defer s.Close()
	err = s.Start()
	assert.NoError(t, err)
	c, err := api.NewClient(http.ClientConfig{Address: "tcp://" + s.Addr, Timeout: time.Minute, KeepAlive: time.Minute})
	assert.NoError(t, err)
	assert.NotNil(t, c)
	p, err := c.GetPortAvailable("127.0.0.1")
	assert.NoError(t, err)
	assert.NotZero(t, p)
	err = c.StartModule(&config.Module{Config: module.Config{Name: "name"}})
	assert.NoError(t, err)
	err = c.StopModule("name")
	assert.NoError(t, err)
}

func TestAPIHttpUnauthorized(t *testing.T) {
	s, err := api.NewServer(&mockEngine{pass: false}, http.ServerConfig{Address: "tcp://127.0.0.1:0", Timeout: time.Minute})
	assert.NoError(t, err)
	defer s.Close()
	err = s.Start()
	assert.NoError(t, err)
	c, err := api.NewClient(http.ClientConfig{Address: "tcp://" + s.Addr, Timeout: time.Minute, KeepAlive: time.Minute, Username: "test"})
	assert.NoError(t, err)
	assert.NotNil(t, c)
	_, err = c.GetPortAvailable("127.0.0.1")
	assert.EqualError(t, err, "[400] Account (test) unauthorized")
	err = c.StartModule(&config.Module{Config: module.Config{Name: "name"}})
	assert.EqualError(t, err, "[400] Account (test) unauthorized")
	err = c.StopModule("name")
	assert.EqualError(t, err, "[400] Account (test) unauthorized")
}
