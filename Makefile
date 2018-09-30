clean:
	@rm -rf temp


# See https://grpc.io/docs/quickstart/go.html
# protoc-gen-go should be added to PATH via export PATH=$PATH:$GOPATH/bin
#
install-deps:
	@apt install tesseract-ocr libtesseract-dev -y
	@go get -u google.golang.org/grpc
	@go get -u github.com/golang/protobuf/protoc-gen-go

install-protoc-linux:
	wget -O ~/Downloads/protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip
	unzip ~/Downloads/protoc.zip -d ~/Downloads/protoc
	sudo mv ~/Downloads/protoc/bin/* /usr/bin/
	chmod +wrx /usr/bin/protoc
	rm -rf ~/Downloads/protoc

gen-protoc:
	protoc -I=./ ./grpc/service.proto --go_out=plugins=grpc:./

build:
	go get github.com/otiai10/gosseract

setup:
	mkdir -p temp

run:
	go run main.go `pwd`/samples/pdf-test.pdf

docker-build:
	docker build

docker-run:
	docker run -d ocr-api

grpc-client:
	go run client/grpcClient.go

grpc-server:
	go run server/grpcServer.go --port 80