package backends_test

import (
	"os"
	"testing"

	"github.com/anzersy/machinery/v2/backends"
	"github.com/anzersy/machinery/v2/config"
	"github.com/anzersy/machinery/v2/tasks"
	"github.com/stretchr/testify/assert"
)

var (
	groupUUID = "123456"
	taskUUIDs = []string{"1", "2", "3"}
)

func initTestMongodbBackend() (backends.Interface, error) {
	cnf := &config.Config{
		ResultBackend:   os.Getenv("MONGODB_URL"),
		ResultsExpireIn: 30,
	}
	backend := backends.NewMongodbBackend(cnf)

	backend.PurgeGroupMeta(groupUUID)
	for _, taskUUID := range taskUUIDs {
		backend.PurgeState(taskUUID)
	}

	if err := backend.InitGroup(groupUUID, taskUUIDs); err != nil {
		return nil, err
	}
	return backend, nil
}

func TestNewMongodbBackend(t *testing.T) {
	if os.Getenv("MONGODB_URL") == "" {
		return
	}

	backend, err := initTestMongodbBackend()
	if assert.NoError(t, err) {
		assert.NotNil(t, backend)
	}
}

func TestSetStatePending(t *testing.T) {
	if os.Getenv("MONGODB_URL") == "" {
		return
	}

	backend, err := initTestMongodbBackend()
	if err != nil {
		t.Fatal(err)
	}

	err = backend.SetStatePending(&tasks.Signature{
		UUID: taskUUIDs[0],
	})
	if assert.NoError(t, err) {
		taskState, err := backend.GetState(taskUUIDs[0])
		if assert.NoError(t, err) {
			assert.Equal(t, tasks.StatePending, taskState.State, "Not StatePending")
		}
	}
}

func TestSetStateReceived(t *testing.T) {
	if os.Getenv("MONGODB_URL") == "" {
		return
	}

	backend, err := initTestMongodbBackend()
	if err != nil {
		t.Fatal(err)
	}

	err = backend.SetStateReceived(&tasks.Signature{
		UUID: taskUUIDs[0],
	})
	if assert.NoError(t, err) {
		taskState, err := backend.GetState(taskUUIDs[0])
		if assert.NoError(t, err) {
			assert.Equal(t, tasks.StateReceived, taskState.State, "Not StateReceived")
		}
	}
}

func TestSetStateStarted(t *testing.T) {
	if os.Getenv("MONGODB_URL") == "" {
		return
	}

	backend, err := initTestMongodbBackend()
	if err != nil {
		t.Fatal(err)
	}

	err = backend.SetStateStarted(&tasks.Signature{
		UUID: taskUUIDs[0],
	})
	if assert.NoError(t, err) {
		taskState, err := backend.GetState(taskUUIDs[0])
		if assert.NoError(t, err) {
			assert.Equal(t, tasks.StateStarted, taskState.State, "Not StateStarted")
		}
	}
}

func TestSetStateSuccess(t *testing.T) {
	if os.Getenv("MONGODB_URL") == "" {
		return
	}

	resultType := "float64"
	resultValue := float64(88.5)

	backend, err := initTestMongodbBackend()
	if err != nil {
		t.Fatal(err)
	}

	signature := &tasks.Signature{
		UUID: taskUUIDs[0],
	}
	taskResults := []*tasks.TaskResult{
		{
			Type:  resultType,
			Value: resultValue,
		},
	}
	err = backend.SetStateSuccess(signature, taskResults)
	assert.NoError(t, err)

	taskState, err := backend.GetState(taskUUIDs[0])
	assert.NoError(t, err)
	assert.Equal(t, tasks.StateSuccess, taskState.State, "Not StateSuccess")
	assert.Equal(t, resultType, taskState.Results[0].Type, "Wrong result type")
	assert.Equal(t, float64(resultValue), taskState.Results[0].Value.(float64), "Wrong result value")
}

func TestSetStateFailure(t *testing.T) {
	if os.Getenv("MONGODB_URL") == "" {
		return
	}

	failStrig := "Fail is ok"

	backend, err := initTestMongodbBackend()
	if err != nil {
		t.Fatal(err)
	}

	signature := &tasks.Signature{
		UUID: taskUUIDs[0],
	}
	err = backend.SetStateFailure(signature, failStrig)
	assert.NoError(t, err)

	taskState, err := backend.GetState(taskUUIDs[0])
	assert.NoError(t, err)
	assert.Equal(t, tasks.StateFailure, taskState.State, "Not StateSuccess")
	assert.Equal(t, failStrig, taskState.Error, "Wrong fail error")
}

func TestGroupCompleted(t *testing.T) {
	if os.Getenv("MONGODB_URL") == "" {
		return
	}

	backend, err := initTestMongodbBackend()
	if err != nil {
		t.Fatal(err)
	}
	taskResultsState := make(map[string]string)

	isCompleted, err := backend.GroupCompleted(groupUUID, len(taskUUIDs))
	if assert.NoError(t, err) {
		assert.False(t, isCompleted, "Actually group is not completed")
	}

	signature := &tasks.Signature{
		UUID: taskUUIDs[0],
	}
	err = backend.SetStateFailure(signature, "Fail is ok")
	assert.NoError(t, err)
	taskResultsState[taskUUIDs[0]] = tasks.StateFailure

	signature = &tasks.Signature{
		UUID: taskUUIDs[1],
	}
	taskResults := []*tasks.TaskResult{
		{
			Type:  "string",
			Value: "Result ok",
		},
	}
	err = backend.SetStateSuccess(signature, taskResults)
	assert.NoError(t, err)
	taskResultsState[taskUUIDs[1]] = tasks.StateSuccess

	signature = &tasks.Signature{
		UUID: taskUUIDs[2],
	}
	err = backend.SetStateSuccess(signature, taskResults)
	assert.NoError(t, err)
	taskResultsState[taskUUIDs[2]] = tasks.StateSuccess

	isCompleted, err = backend.GroupCompleted(groupUUID, len(taskUUIDs))
	if assert.NoError(t, err) {
		assert.True(t, isCompleted, "Actually group is completed")
	}

	taskStates, err := backend.GroupTaskStates(groupUUID, len(taskUUIDs))
	assert.NoError(t, err)

	assert.Equal(t, len(taskStates), len(taskUUIDs), "Wrong len tasksStates")
	for _, taskState := range taskStates {
		assert.Equal(
			t,
			taskResultsState[taskState.TaskUUID],
			taskState.State,
			"Wrong state on", taskState.TaskUUID,
		)
	}
}

func TestGroupStates(t *testing.T) {
	if os.Getenv("MONGODB_URL") == "" {
		return
	}

	backend, err := initTestMongodbBackend()
	if err != nil {
		t.Fatal(err)
	}

	taskStates, err := backend.GroupTaskStates(groupUUID, len(taskUUIDs))
	assert.NoError(t, err)
	for i, taskState := range taskStates {
		assert.Equal(t, taskUUIDs[i], taskState.TaskUUID)
	}
}
