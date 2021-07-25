SHELL=/bin/bash
#!make
include envfile
export $(shell sed 's/=.*//' envfile)

# to see all colors, run
# bash -c 'for c in {0..255}; do tput setaf $c; tput setaf $c | cat -v; echo =$c; done'
# the first 15 entries are the 8-bit colors

# define standard colors
ifneq (,$(findstring xterm,${TERM}))
	BLACK        := $(shell tput -Txterm setaf 0)
	RED          := $(shell tput -Txterm setaf 1)
	GREEN        := $(shell tput -Txterm setaf 2)
	YELLOW       := $(shell tput -Txterm setaf 3)
	LIGHTPURPLE  := $(shell tput -Txterm setaf 4)
	PURPLE       := $(shell tput -Txterm setaf 5)
	BLUE         := $(shell tput -Txterm setaf 6)
	WHITE        := $(shell tput -Txterm setaf 7)
	RESET := $(shell tput -Txterm sgr0)
else
	BLACK        := ""
	RED          := ""
	GREEN        := ""
	YELLOW       := ""
	LIGHTPURPLE  := ""
	PURPLE       := ""
	BLUE         := ""
	WHITE        := ""
	RESET        := ""
endif

# set target color
TARGET_COLOR := $(BLUE)

POUND = \#

SLACKIT_VERSION = 0.0.1
CURLIT_VERSION = 0.0.1


.PHONY: no_targets__ info help build deploy doc
	no_targets__: help

.DEFAULT_GOAL := help

colors: ## show all the colors
	@echo "${BLACK}BLACK${RESET}"
	@echo "${RED}RED${RESET}"
	@echo "${GREEN}GREEN${RESET}"
	@echo "${YELLOW}YELLOW${RESET}"
	@echo "${LIGHTPURPLE}LIGHTPURPLE${RESET}"
	@echo "${PURPLE}PURPLE${RESET}"
	@echo "${BLUE}BLUE${RESET}"
	@echo "${WHITE}WHITE${RESET}"
	
build:  ## run the docker builds
	@echo "${GREEN}Running build.${RESET}"
    golint
	docker build -t leewallen/slackit:${SLACKIT_VERSION} -f Dockerfile .
	docker build -t leewallen/curlit:${CURLIT_VERSION} -f ScheduledDockerfile .
	@echo "${GREEN}Build finished.${RESET}"

publish:  ## push images to Docker repository
	@echo "${GREEN}Publishing docker images.${RESET}"
	docker push leewallen/slackit:${SLACKIT_VERSION}
	docker push leewallen/curlit:${CURLIT_VERSION}
	@echo "${GREEN}Publish finished.${RESET}"

deploy:  ## deploy kubernetes related resources
	@echo "${GREEN}Deploying to k8s cluster.${RESET}"
	env
	kubectl create configmap slackit-configmap --from-literal=SLACK_URL="${SLACK_URL}"
	kubectl create configmap swanson-configmap --from-literal=SWANSON_CHANNEL="${SWANSON_CHANNEL}" --from-literal=SWANSON_URL="${SWANSON_URL}"
	kubectl create configmap nasa-configmap --from-literal=NASA_CHANNEL="${NASA_CHANNEL}" --from-literal=NASA_URL="${NASA_URL}"
	kubectl create configmap xkcd-configmap --from-literal=XKCD_CHANNEL="${XKCD_CHANNEL}" --from-literal=XKCD_URL="${XKCD_URL}"

	kubectl apply -f ./k8s/schedule-nasa.yaml
	kubectl apply -f ./k8s/schedule-xkcd.yaml
	kubectl apply -f ./k8s/schedule-swanson.yaml
	kubectl apply -f ./k8s/service-slackit.yaml
	kubectl apply -f ./k8s/deployment-slackit.yaml
	@echo "${GREEN}deploy finished${RESET}"

delete: ## delete kubernetes related resources
	@echo "${LIGHTPURPLE}Delete deployment.${RESET}"
	kubectl delete configmap slackit-configmap
	kubectl delete configmap swanson-configmap
	kubectl delete configmap nasa-configmap
	kubectl delete configmap xkcd-configmap
	kubectl delete -f ./k8s/service-slackit.yaml
	kubectl delete -f ./k8s/schedule-swanson.yaml
	kubectl delete -f ./k8s/schedule-nasa.yaml
	kubectl delete -f ./k8s/schedule-xkcd.yaml
	kubectl delete -f ./k8s/deployment-slackit.yaml
	@echo "${LIGHTPURPLE}deployment deleted${RESET}"

describe: ## describe a pod that has the app set to get-swanson-quote
	kubectl describe pod --selector=app=slackit

pods: ## get pods that have the app set to get-swanson-quote
	kubectl get pods --selector=app=slackit

help:
	@echo "${BLACK}-----------------------------------------------------------------${RESET}"
	@grep -E '^[a-zA-Z_0-9%-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "${TARGET_COLOR}%-30s${RESET} %s\n", $$1, $$2}'

