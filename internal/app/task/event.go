package task

import "time"

// TaskCreated is the fact announced on the broker whenever a task comes into
// existence — created through the HTTP API or as a payment reminder alike.
// It is a flat, self-contained fact stating what happened and naming no
// action; consumers own their reactions. The JSON tags are the cross-service
// payload contract:
// the owner's id is numeric to match the billing contract's user ids, a
// conversion from the internal string representation made here, at the event
// boundary, mirroring the inbound conversion.
type TaskCreated struct {
	TaskId      int       `json:"taskId"`
	UserId      int       `json:"userId"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Publisher is the seam for announcing task facts on the broker. The task
// package defines it narrow: publication is best-effort, so the
// service can swallow a failure and still report the save's result.
type Publisher interface {
	PublishTaskCreated(event TaskCreated) error
}
