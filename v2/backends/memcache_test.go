package backends_test

import (
	"os"
	"testing"
	"time"

	"github.com/anzersy/machinery/v2/backends"
	"github.com/anzersy/machinery/v2/config"
	"github.com/anzersy/machinery/v2/tasks"
	"github.com/stretchr/testify/assert"
)

func TestGroupCompletedMemcache(t *testing.T) {
	memcacheURL := os.Getenv("MEMCACHE_URL")
	if memcacheURL == "" {
		return
	}

	groupUUID := "testGroupUUID"
	task1 := &tasks.Signature{
		UUID:      "testTaskUUID1",
		GroupUUID: groupUUID,
	}
	task2 := &tasks.Signature{
		UUID:      "testTaskUUID2",
		GroupUUID: groupUUID,
	}

	backend := backends.NewMemcacheBackend(new(config.Config), []string{memcacheURL})

	// Cleanup before the test
	backend.PurgeState(task1.UUID)
	backend.PurgeState(task2.UUID)
	backend.PurgeGroupMeta(groupUUID)

	groupCompleted, err := backend.GroupCompleted(groupUUID, 2)
	if assert.Error(t, err) {
		assert.False(t, groupCompleted)
		assert.Equal(t, "memcache: cache miss", err.Error())
	}

	backend.InitGroup(groupUUID, []string{task1.UUID, task2.UUID})

	groupCompleted, err = backend.GroupCompleted(groupUUID, 2)
	if assert.Error(t, err) {
		assert.False(t, groupCompleted)
		assert.Equal(t, "memcache: cache miss", err.Error())
	}

	backend.SetStatePending(task1)
	backend.SetStateStarted(task2)
	groupCompleted, err = backend.GroupCompleted(groupUUID, 2)
	if assert.NoError(t, err) {
		assert.False(t, groupCompleted)
	}

	taskResults := []*tasks.TaskResult{new(tasks.TaskResult)}
	backend.SetStateStarted(task1)
	backend.SetStateSuccess(task2, taskResults)
	groupCompleted, err = backend.GroupCompleted(groupUUID, 2)
	if assert.NoError(t, err) {
		assert.False(t, groupCompleted)
	}

	backend.SetStateFailure(task1, "Some error")
	groupCompleted, err = backend.GroupCompleted(groupUUID, 2)
	if assert.NoError(t, err) {
		assert.True(t, groupCompleted)
	}
}

func TestGetStateMemcache(t *testing.T) {
	memcacheURL := os.Getenv("MEMCACHE_URL")
	if memcacheURL == "" {
		return
	}

	signature := &tasks.Signature{
		UUID:      "testTaskUUID",
		GroupUUID: "testGroupUUID",
	}

	backend := backends.NewMemcacheBackend(new(config.Config), []string{memcacheURL})

	go func() {
		backend.SetStatePending(signature)
		<-time.After(2 * time.Millisecond)
		backend.SetStateReceived(signature)
		<-time.After(2 * time.Millisecond)
		backend.SetStateStarted(signature)
		<-time.After(2 * time.Millisecond)
		taskResults := []*tasks.TaskResult{
			{
				Type:  "float64",
				Value: 2,
			},
		}
		backend.SetStateSuccess(signature, taskResults)
	}()

	var (
		taskState *tasks.TaskState
		err       error
	)
	for {
		taskState, err = backend.GetState(signature.UUID)
		if taskState == nil {
			assert.Equal(t, "memcache: cache miss", err.Error())
			continue
		}

		assert.NoError(t, err)
		if taskState.IsCompleted() {
			break
		}
	}
}

func TestPurgeStateMemcache(t *testing.T) {
	memcacheURL := os.Getenv("MEMCACHE_URL")
	if memcacheURL == "" {
		return
	}

	signature := &tasks.Signature{
		UUID:      "testTaskUUID",
		GroupUUID: "testGroupUUID",
	}

	backend := backends.NewMemcacheBackend(new(config.Config), []string{memcacheURL})

	backend.SetStatePending(signature)
	taskState, err := backend.GetState(signature.UUID)
	assert.NotNil(t, taskState)
	assert.NoError(t, err)

	backend.PurgeState(taskState.TaskUUID)
	taskState, err = backend.GetState(signature.UUID)
	assert.Nil(t, taskState)
	assert.Error(t, err)
}
