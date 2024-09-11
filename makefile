build: build-image build-discovery build-gateway build-analytics
	@echo "Building all"
build-image:
	@echo "Building image"
	docker build -t image-service ./image_service
build-analytics:
	@echo "Building analytics"
	docker build -t analytics-service ./recommendation_service
build-gateway:
	@echo "Building gateway"
	docker build -t gateway-service ./gateway
build-discovery:
	@echo "Building discovery"
	docker build -t discovery-service ./service_discovery

push: push-image push-discovery push-gateway push-analytics
	@echo "Pushing all"
push-image:
	@echo "Pushing image"
	docker tag image-service:latest stfnbcx/image-service:latest
	docker push stfnbcx/image-service:latest
push-analytics:	
	docker tag analytics-service:latest stfnbcx/analytics-service:latest
	docker push stfnbcx/analytics-service:latest
push-gateway:
	docker tag gateway-service:latest stfnbcx/gateway-service:latest
	docker push stfnbcx/gateway-service:latest
push-discovery:
	docker tag discovery-service:latest stfnbcx/discovery-service:latest
	docker push stfnbcx/discovery-service:latest

pull: pull-image pull-discovery pull-gateway pull-analytics
	@echo "Pulling all"
pull-image:
	@echo "Pulling image"
	docker pull stfnbcx/image-service:latest
pull-analytics:
	docker pull stfnbcx/analytics-service:latest
pull-gateway:
	docker pull stfnbcx/gateway-service:latest
pull-discovery:
	docker pull stfnbcx/discovery-service:latest