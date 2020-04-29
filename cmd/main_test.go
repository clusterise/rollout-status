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

func assertRolloutFailure(t *testing.T, expectedMessage string) {
	wrapper := mockWrapperFromAssets(t.Name())
	rolloutStatus := status.TestRollout(wrapper, IgnoredByMock, IgnoredByMock)

	re, ok := rolloutStatus.Error.(status.RolloutError)
	assert.True(t, ok)

	assert.False(t, rolloutStatus.Continue)
	assert.Equal(t, expectedMessage, re.Message)
}

func TestLimitRange(t *testing.T) {
	assertRolloutFailure(t, `replicaset limit-range-7dfcd777fd failed to create pods: Pod "limit-range-7dfcd777fd-kmgsc" is invalid: spec.containers[0].resources.requests: Invalid value: "200Mi": must be less than or equal to memory limit`)
}

func TestResourceQuota(t *testing.T) {
	assertRolloutFailure(t, `replicaset resource-quota-6884c5558d failed to create pods: pods "resource-quota-6884c5558d-jzwxq" is forbidden: exceeded quota: main, requested: memory=200Mi, used: memory=0, limited: memory=100Mi`)
}

func TestImagePullBackOff(t *testing.T) {
	assertRolloutFailure(t, `container main is in ImagePullBackOff: Back-off pulling image "bogus-image:does-not-exist"`)
}

func TestErrImagePull(t *testing.T) {
	assertRolloutFailure(t, `container main is in ErrImagePull: rpc error: code = Unknown desc = Error response from daemon: pull access denied for bogus-image, repository does not exist or may require 'docker login': denied: requested access to the resource is denied`)
}

func TestPending(t *testing.T) {
	assertRolloutFailure(t, `failed to scheduled pod: 0/1 nodes are available: 1 Insufficient memory.`)
}

func TestRunContainerError(t *testing.T) {
	assertRolloutFailure(t, `container main is in RunContainerError: failed to start container "7e28d6efbe847126d763f15d88a68a225bb9b746651ca82abeba419e315e7c18": Error response from daemon: OCI runtime create failed: container_linux.go:349: starting container process caused "exec: \"binary-not-found\": executable file not found in $PATH": unknown`)
}

func TestCrashloopBackoff(t *testing.T) {
	assertRolloutFailure(t, `container main is in CrashLoopBackOff: back-off 10s restarting failed container=main pod=crashloop-backoff-58894954f8-68pcr_default(6c016294-6f4c-4f68-a919-3f8385dee8fd)`)
}

func TestReadiness(t *testing.T) {
	assertRolloutFailure(t, `deployment readiness is not progressing: ReplicaSet "readiness-78d744cf44" has timed out progressing.`)
}
