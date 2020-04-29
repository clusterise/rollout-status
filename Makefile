REPO = rollout-status
TAG = 1.0

all: build

build: Dockerfile preflight
	docker build -t $(REPO):$(TAG) .

preflight:
	go mod vendor
	go fmt dite.pro/rollout-status/...

.PHONY: test
test:
	./test-e2e.sh $(REPO):$(TAG)

.PHONY: cleanup
cleanup:
	rm -rf vendor
	minikube stop
	minikube delete
