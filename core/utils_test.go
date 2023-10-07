package core

import "time"

type MockTimer struct {
	Date time.Time
}

func (m *MockTimer) Get() time.Time {
	return m.Date
}

func createMockTimer(constDate time.Time) *MockTimer {
	return &MockTimer{Date: constDate}
}
