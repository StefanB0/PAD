build-all: build-image build-discovery build-gateway build-analytics
	@echo "Building all"
build-image:
	@echo "Building image"
	docker build -t image-service ./image_service
	docker tag image-service:latest stfnbcx/image-service:latest
	docker push stfnbcx/image-service:latest
build-analytics:
	@echo "Building analytics"
	docker build -t analytics-service ./recommendation_service
	docker tag analytics-service:latest stfnbcx/analytics-service:latest
	docker push stfnbcx/analytics-service:latest
build-gateway:
	@echo "Building gateway"
	docker build -t gateway-service ./gateway
	docker tag gateway-service:latest stfnbcx/gateway-service:latest
	docker push stfnbcx/gateway-service:latest
build-discovery:
	@echo "Building discovery"
	docker build -t discovery-service ./service_discovery
	docker tag discovery-service:latest stfnbcx/discovery-service:latest
	docker push stfnbcx/discovery-service:latest
	
