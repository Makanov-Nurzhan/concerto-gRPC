PROTO_DIR = api/proto
PROTO_DIR = api/proto
GEN_DIR   = api/gen/adminv1

PROTO_FILES = $(PROTO_DIR)/concerto_admin.proto

proto: $(PROTO_FILES)
	protoc \
		-I=$(PROTO_DIR) \
		--go_out=$(GEN_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)


install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

proto-all: install-tools proto
