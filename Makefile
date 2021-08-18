## 
GO111MODULE=on
## image registry host
REGISTRY?=asia.gcr.io
#
PROJECTID?=second-hand-boutique
## git commit
COMMIT?=$(shell git rev-parse HEAD)
## build date
DATE?=$(shell date "+%Y/%m/%dT%H:%M:%S")
## app name
APP_NAME?=$(shell basename ${PWD})
## image name
IMAGE_NAME?=${REGISTRY}/${PROJECTID}/${APP_NAME}:${TAG}
## server port
PORT?=58867
## server consul
CONSUL?=localhost:8500

#dev-deploy
DEPLOY_API_NAME:=$(shell echo ${IMAGE_NAME} | sed -e 's/\//\\\//g')
OLD_API_NAME:=$(shell cat ${PWD}/deploy.yaml |grep harbor |awk '{print $$2}'|sed -e 's/\//\\\//g')

## build: go build the application
.PHONY: build
build: clean
	go build -o api .

## run: run the application server arg: PORT , CONSUL
.PHONY: run
run:
	go run main.go server -p ${PORT} -c ${CONSUL}

## clean: golang clena
.PHONY: clean
clean:
	go clean
	docker rm 

## build-for-docker: use local machine to build the service and tar into the image
.PHONY: build-for-docker
build-for-docker: clean check-tag
	@echo TAG ${TAG}

	@echo "go Build"
	CGO_ENABLED=1 GOOS=linux go build  -a -installsuffix cgo \
    -ldflags "-X main.VERSION=${TAG} -X main.COMMIT=${COMMIT} -X main.BUILD=${DATE}" \
    -o api

	@echo "docker build image ${IMAGE_NAME}"
	docker build -t ${IMAGE_NAME} .

## build-in-docker: use docker to build the service and tar into the image
.PHONY: build-in-docker
build-in-docker:

	@echo "docker build"
	docker build . \
                  --build-arg COMMIT=${COMMIT} \
                  --build-arg DATE=${DATE} \
                  --build-arg TAG=${TAG} \
                  --file Dockerfile --tag  ${IMAGE_NAME}

	docker image prune -f --filter label=stage=builder


check-tag:
ifndef TAG
	$(error TAG not set)
endif

# change deploy.yaml image version
.PHONY: update-deploy
update-deploy:
	@sed -i 's/${OLD_API_NAME}/${DEPLOY_API_NAME}/g' ${PWD}/deploy.yaml

## help: help range
.PHONY: help
help:
	@echo "Useage:\n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

