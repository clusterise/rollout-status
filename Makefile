REPO = rollout-status
TAG = 1.0

all: build

build: Dockerfile
	docker build -t $(REPO):$(TAG) .
