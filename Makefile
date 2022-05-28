.PHONY: protos

protos:
	 protoc  --proto_path=protos --go_out=protos --go-grpc_out=require_unimplemented_servers=false:protos currency.proto