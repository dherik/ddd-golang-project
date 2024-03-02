package integration

import (
	"testing"

	"github.com/dherik/ddd-golang-project/tests/integration/setup"
	"github.com/stretchr/testify/suite"
)

//TODO login test
//TODO login error test
//TODO try to access protect endpoint not authenticated test
//TODO try to access protect endpoint authenticated test
//TODO create user and login test

type UserTestSuite struct {
	suite.Suite
}

func (s *UserTestSuite) TestLogin() {
	if testing.Short() {
		s.T().Skip("Skip test for postgresql repository")
	}

	token, err := setup.Login("admin", "some_password")

	if err != nil {
		s.T().Fatalf("error while login: %v", err)
	}

	if token == "" {
		s.T().Fatalf("token should not be empty")
	}
}
