COMMON_PATH	?= $(shell pwd)

tidy:
	GOPRIVATE=$(GOPRIVATE) go mod tidy

update:
	GOPRIVATE=$(GOPRIVATE) go get -u ./...

generate-proto:
	cd protocols/delivery-pkg && protoc --go_out=paths=source_relative:. --go_opt=paths=source_relative  --go-grpc_out=paths=source_relative:. --go-grpc_opt=paths=source_relative  *.proto
