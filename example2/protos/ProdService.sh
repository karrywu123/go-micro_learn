
go get -u github.com/asim/go-micro/v3

go get -u github.com/asim/go-micro/plugins/registry/consul/v3

go get -u github.com/golang/protobuf/protoc-gen-go

go get -u github.com/micro/micro/v3/cmd/protoc-gen-micro

protoc -I . --go_out=plugins=grpc:. ./Models.proto

protoc -I . --go_out=plugins=grpc:. ./ProdService.proto


protoc --micro_out=. --go_out=. ./Models.proto

protoc --micro_out=. --go_out=. ./ProdService.proto

protoc-go-inject-tag -input=. ./services/Models.pb.go


