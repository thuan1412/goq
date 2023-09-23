package task

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/samber/lo"
)

const (
	configField    = "Config"
	taskNameField  = "Name"
	taskQueueField = "Queue"
	defaultQueue   = "default"
)

var (
	invalidTaskerErr     = errors.New("invalid tasker object")
	invalidTaskConfigErr = errors.New("invalid tasker config object")
)

func GetTaskName(tasker Tasker) string {
	obj := reflect.Indirect(reflect.ValueOf(tasker))
	config := obj.FieldByName(configField)
	if !config.IsValid() {
		panic(invalidTaskerErr)
	}
	if config.IsZero() {
		return fmt.Sprintf("%T", tasker)
	}

	nameField := config.FieldByName(taskNameField)
	if !nameField.IsValid() {
		panic(invalidTaskerErr)
	}
	taskName := nameField.String()

	return taskName
}

func GetTaskQueue(tasker Tasker) string {
	obj := reflect.Indirect(reflect.ValueOf(tasker))
	config := obj.FieldByName(configField)
	if !config.IsValid() {
		panic(invalidTaskerErr)
	}

	queueField := config.FieldByName(taskQueueField)
	if !queueField.IsValid() {
		panic(invalidTaskConfigErr)
	}
	queueName := queueField.String()

	return lo.Ternary(queueName == "", defaultQueue, queueName)
}
