package cron

import "context"

// -------------------------------------------------------------------------------
// serial basic task
// -------------------------------------------------------------------------------

// NewSerialTask creates a task that run only serially, provided an
// action and a policy
func NewSerialTask(action TaskAction) Task {
	name := generateTaskName()
	return &basicSerialTask{name: name, action: action}
}

// NewSerialTaskWithName creates a task that can only be run serially given an
// action, name, and policy
func NewSerialTaskWithName(name string, action TaskAction) Task {
	return &basicSerialTask{name: name, action: action}
}

type basicSerialTask struct {
	name   string
	action TaskAction
}

// Name returns the name of a basic serial task
func (bst basicSerialTask) Name() string {
	return bst.name
}

// Execute runs the action that was assigned for the task
func (bst basicSerialTask) Execute(ctx context.Context) error {
	return bst.action(ctx)
}

func (bst basicSerialTask) OnStart()             {}
func (bst basicSerialTask) OnCancellation()      {}
func (bst basicSerialTask) OnComplete(err error) {}
func (bst basicSerialTask) Serial()              {}
