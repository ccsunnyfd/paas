API_PROTO_FILES=$(shell find api -name *.proto)

.PHONY: grpc
# generate grpc code
grpc:
	protoc --proto_path=. \
           --go_out=paths=source_relative:. \
           --micro_out=paths=source_relative:. \
           $(API_PROTO_FILES)
