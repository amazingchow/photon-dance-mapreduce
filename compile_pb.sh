#!/bin/bash
# pb version: 3.3.0

if [ -z "$GOMODULEPATH" ]; then
    echo "YOU MUST SET ENV GOMODULEPATH."
fi

rm -rf api

PROTO_INSTALL=/usr/local

cd $GOMODULEPATH
for i in $(ls $GOMODULEPATH/github.com/amazingchow/mapreduce/pb/*.proto); do
	fn=github.com/amazingchow/mapreduce/pb/$(basename "$i")
	echo "compile" $fn
	$PROTO_INSTALL/bin/protoc -I$PROTO_INSTALL/include -I . \
		-I$GOMODULEPATH \
		-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:. "$fn"
	$PROTO_INSTALL/bin/protoc -I$PROTO_INSTALL/include -I . \
		-I$GOMODULEPATH \
		-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--grpc-gateway_out=logtostderr=true:. "$fn"
	$PROTO_INSTALL/bin/protoc -I$PROTO_INSTALL/include -I . \
		-I$GOMODULEPATH \
		-I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--swagger_out=logtostderr=true:. "$fn"
done
