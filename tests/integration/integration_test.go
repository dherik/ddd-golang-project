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
	setup.ResetData()
	setup.LoadDML()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExecuteAllSuites(t *testing.T) {

	if testing.Short() {
		t.Skip("Skip integration tests execution")
	}

	setup.SetupIntegrationTest()

	suite.Run(t, new(TaskTestSuite))
	suite.Run(t, new(UserTestSuite))

	setup.StopIntegrationTest()

}
