protobuild:
	@protoc --go_out="." --go-grpc_out="." ./proto/*.proto

get-deps:
	@go get -u aidanwoods.dev/go-paseto
	@go get github.com/joho/godotenv
	@go get -u github.com/rs/zerolog/log
	@go get go.mongodb.org/mongo-driver/mongo
	@go get -u gorm.io/gorm
	@go get -u gorm.io/driver/postgres