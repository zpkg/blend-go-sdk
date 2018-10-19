package cron

import (
	"context"
	"fmt"

	"github.com/blend/go-sdk/uuid"
)

// --------------------------------------------------------------------------------
// quick task creation
// --------------------------------------------------------------------------------

// NewTask returns a new task wrapper for a given TaskAction.
func NewTask(action TaskAction) Task {
	name := generateTaskName()
	return &basicTask{name: name, action: action}
}

// NewTaskWithName returns a new task wrapper with a given name for a given TaskAction.
func NewTaskWithName(name string, action TaskAction) Task {
	return &basicTask{name: name, action: action}
}

type basicTask struct {
	name   string
	action TaskAction
}

func (bt basicTask) Name() string {
	return bt.name
}
func (bt basicTask) Execute(ctx context.Context) error {
	return bt.action(ctx)
}
func (bt basicTask) OnStart()             {}
func (bt basicTask) OnCancellation()      {}
func (bt basicTask) OnComplete(err error) {}

// generateTaskName returns a unique identifier that can be used to name/tag tasks
func generateTaskName() string {
	return fmt.Sprintf("task_%s", uuid.V4().ToShortString())
}
