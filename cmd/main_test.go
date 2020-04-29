package main_test

import (
	"dite.pro/rollout-status/pkg/status"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func TestNotFound(t *testing.T) {
	wrapper := mockWrapper([]appsv1.Deployment{}, []appsv1.ReplicaSet{}, []v1.Pod{})
	rolloutStatus := status.TestRollout(wrapper, "ns-alpha", "foo=bar")

	re, ok := rolloutStatus.Error.(status.RolloutError)
	assert.True(t, ok)

	assert.False(t, rolloutStatus.Continue)
	assert.Equal(t, `Selector "foo=bar" did not match any deployments in namespace "ns-alpha"`, re.Message)
}

func TestSuccess(t *testing.T) {
	wrapper := mockWrapperFromAssets(t.Name())
	rolloutStatus := status.TestRollout(wrapper, "ns-alpha", "foo=bar")

	assert.False(t, rolloutStatus.Continue)
	assert.Nil(t, rolloutStatus.Error)
}
