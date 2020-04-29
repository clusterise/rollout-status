Rollout Status
==============

An alternative to `kubectl rollout status`.

It is intended as a building block for your CI/CD pipeline. Unlike `kubectl` it does not rely *only* on `progressDeadlineSeconds` and fails fast when appropriate, such as when image does not exist or environment properties are misconfigured. 


Usage
-----

Docker Image is [available on Docker Hub](https://hub.docker.com/repository/docker/clusterise/rollout-status
).

```console
docker run -v "$HOME/.kube:/root/.kube:ro" clusterise/rollout-status:1.0 -namespace "$NAMESPACE" -selector "$SELECTOR"
```

Since the script is intended for a CI/CI, you will likely copy the binary to your own `linux-amd64` image:

```Dockerfile
# ...

COPY --from=clusterise/rollout-status:1.0 /rollout-status /opt/rollout-status 
```

Example
-------

```console
go run github.com/clusterise/rollout-status/cmd -selector app=crashloop-backoff
```
```json
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


Output
------

See `schema.json` for details. The program blocks until rollout succeeds (as `success=true`) or is not progressing and effectively failed (`success=false`). In the latter case an `error` object is returned. It contains all details that should be neccessary to remediate the issue. Besides `error.mesage` your pipeline should also switch on `error.code` and print help specific to your pipeline.
