package db

import (
	"github.com/stretchr/testify/mock"
	"lehuuthoct/lht-microservices/leader/model"
)

// init mock object for bolt db client
// instead of using a bolt.DB pointer
type MockBoltClient struct {
	mock.Mock
}

// Mock 3 methods to fulfill IBoltClient interface
func (m *MockBoltClient) FindLeaderById(leaderId string) (model.Leader, error)  {
	args := m.Mock.Called(leaderId)
	return args.Get(0).(model.Leader), args.Error(1)
}

func (m *MockBoltClient) OpenBoltDB()  {
}

func (m *MockBoltClient) MockLeaderData() {
}

func (m *MockBoltClient) Check() bool  {
	args := m.Mock.Called()
	return args.Get(0).(bool)
}

