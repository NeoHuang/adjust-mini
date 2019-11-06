VERSION = $(shell git rev-parse --short HEAD)
DEPLOY_SERVICE = $(subst _,-,$(service))

.PHONY: all
all: build tag push upgrade

.PHONY: build
build:
	@printf "\033[32mBuild $(service) Docker image\n\033[0m"
	docker build  --build-arg backend_target=$(service) --build-arg version=$(VERSION) ./ -t neohuang/$(service)

.PHONY: tag
tag:
	@printf "\033[32mTag $(service) Docker image\n\033[0m"
	docker tag neohuang/$(service) neohuang/$(service):$(VERSION)

.PHONY: push
push:
	@printf "\033[32mPush $(service) Docker image\n\033[0m"
	docker push neohuang/$(service):$(VERSION)

.PHONY: deploy
deploy:
	@printf "\033[32mDeploy $(service) to k8s cluster\n\033[0m"
	kubectl set image deployment $(DEPLOY_SERVICE)-deployment $(DEPLOY_SERVICE)=neohuang/$(service):$(VERSION)
	@kubectl rollout status deployment $(DEPLOY_SERVICE)-deployment

.PHONY: install
install:
	@printf "\033[32mDeploy $(service) to k8s cluster via helm\n\033[0m"
	helm install --name $(DEPLOY_SERVICE) ./$(service)/$(service)-helm/ --set image.repository=neohuang/$(service),image.tag=$(VERSION)


.PHONY: upgrade
upgrade:
	@printf "\033[32mDeploy $(service) to k8s cluster via helm\n\033[0m"
	helm upgrade $(DEPLOY_SERVICE) ./$(service)/$(service)-helm/ --set image.repository=neohuang/$(service),image.tag=$(VERSION)

.PHONY: restart
restart:
	@printf "\033[32mRolling restart $(service)\n\033[0m"
	@kubectl patch deployment $(DEPLOY_SERVICE)-deployment -p "{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"date\":\"`date +'%s'`\"}}}}}"
	@kubectl rollout status deployment $(DEPLOY_SERVICE)-deployment

