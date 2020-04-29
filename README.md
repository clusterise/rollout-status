Rollout Status
==============

An alternative to `kubectl rollout status`.

It is intended as a building block for your CI/CD pipeline. Unlike `kubectl` it does not rely *only* on `progressDeadlineSeconds` and fails fast when appropriate, such as when image does not exist or environment properties are misconfigured. 

Example
-------

```console
go run dite.pro/rollout-status/cmd -selector app=crashloop-backoff
{
  "success": false,
  "error": {
    "code": "process-crashing",
    "message": "Container \"main\" is in \"CrashLoopBackOff\": back-off 5m0s restarting failed container=main pod=crashloop-backoff-7fd845849c-cfvqd_default(617f6364-5bd4-4e69-b19f-fcf3ce4c171a)",
    "type": "rollout",
    "log": "panic: invalid configuration: [unable to read client-key /Users/mikulas/.minikube/profiles/minikube/client.key for minikube due to open /Users/mikulas/.minikube/profiles/minikube/client.key: no such file or directory, unable to read certificate-authority /Users/mikulas/.minikube/ca.crt for minikube due to open /Users/mikulas/.minikube/ca.crt: no such file or directory]\n\ngoroutine 1 [running]:\nmain.makeClientset(0xc000041440, 0x12, 0x0)\n\t/src/cmd/main.go:57 +0xfb\nmain.main()\n\t/src/cmd/main.go:34 +0x21c\n"
  }
}
```
