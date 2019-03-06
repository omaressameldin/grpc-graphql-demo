FROM znly/protoc:0.4.0
COPY ./proto /proto
RUN mkdir -p /pkg/api/v1
RUN ls /proto/v1

ENTRYPOINT ["protoc", "--proto_path=/proto/v1", "--go_out=plugins=grpc:/pkg/api/v1", "todo-service.proto" ]