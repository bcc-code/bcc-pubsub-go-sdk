package model

import "time"

// DispatchResponse is returned when a task has been scheduled.
type DispatchResponse struct {
	DispatchId   string    `json:"dispatchId"`
	Endpoint     string    `json:"endpoint"`
	ScheduleTime time.Time `json:"scheduleTime"`
}
