package cmd

import (
	"fmt"
	"testing"

	config "bitbucket.org/butenta/pkg-config"
)

type TestConfiguration struct {
	Testing  string
	App      *config.App
	Database *config.Database
}

func (t *TestConfiguration) GetApp() *config.App {
	return t.App
}

func (t *TestConfiguration) GetDB() *config.Database {
	return t.Database
}

func TestRoot(t *testing.T) {
	cmd := New("test", "For test usage Long name", "0.0.1", "test")

	cmd.AddVersion()

	config := &TestConfiguration{}
	cmd.LoadConfig(config)

	cmd.AddService(func() {
		fmt.Println("Service TESTING")
	})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Wanted error nil, got: %s", err)
	}
}
