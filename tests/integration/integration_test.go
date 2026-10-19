package integration

import (
	"log/slog"
	"testing"

	"github.com/dherik/ddd-golang-project/tests/integration/setup"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
}

// Will run after each test in the Suites
func (suite *TaskTestSuite) SetupTest() {
	slog.Info("Reseting database for the next test")
	suite.Require().NoError(setup.ResetData())
	suite.Require().NoError(setup.LoadDML())
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExecuteAllSuites(t *testing.T) {

	if testing.Short() {
		t.Skip("Skip integration tests execution")
	}

	if err := setup.SetupIntegrationTest(); err != nil {
		t.Fatalf("failed to set up integration tests: %v", err)
	}

	suite.Run(t, new(TaskTestSuite))
	suite.Run(t, new(UserTestSuite))

	if err := setup.StopIntegrationTest(); err != nil {
		t.Errorf("failed to tear down integration tests: %v", err)
	}

}
