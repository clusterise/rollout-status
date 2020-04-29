package main_test

import (
	"dite.pro/rollout-status/pkg/status"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"testing"
)

const IgnoredByMock = "any-value"

func TestNotFound(t *testing.T) {
	wrapper := mockWrapper([]appsv1.Deployment{}, []appsv1.ReplicaSet{}, []v1.Pod{})
	rolloutStatus := status.TestRollout(wrapper, "any-ns", "foo=bar")

	re, ok := rolloutStatus.Error.(status.RolloutError)
	assert.True(t, ok)

	assert.False(t, rolloutStatus.Continue)
	assert.Equal(t, `Selector "foo=bar" did not match any deployments in namespace "any-ns"`, re.Message)
}

func TestSuccess(t *testing.T) {
	wrapper := mockWrapperFromAssets(t.Name())
	rolloutStatus := status.TestRollout(wrapper, IgnoredByMock, IgnoredByMock)

	assert.False(t, rolloutStatus.Continue)
	assert.Nil(t, rolloutStatus.Error)
}

func TestLimitRange(t *testing.T) {
	wrapper := mockWrapperFromAssets(t.Name())
	rolloutStatus := status.TestRollout(wrapper, IgnoredByMock, IgnoredByMock)

	re, ok := rolloutStatus.Error.(status.RolloutError)
	assert.True(t, ok)

	assert.False(t, rolloutStatus.Continue)
	assert.Equal(t, `xxx"`, re.Message)
}
