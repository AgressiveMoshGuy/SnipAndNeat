FROM golang:alpine3.19 AS builder

ENV CGO_ENABLED=1

RUN apk add --no-cache gcc musl-dev make

WORKDIR /build
# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .
RUN go mod download -x

COPY . .
RUN make build

FROM alpine:3.19
COPY --from=builder /build/output /bin/output

EXPOSE 8080/tcp
EXPOSE 8081/tcp
EXPOSE 8082/tcp
EXPOSE 8083/tcp
ENTRYPOINT ["/bin/output"]

