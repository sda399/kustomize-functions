
FROM public.ecr.aws/docker/library/golang:1.24-bullseye AS builder
ENV CGO_ENABLED=0
WORKDIR /go/src/
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags '-w -s' -v -o /usr/local/bin/function ./

FROM gcr.io/distroless/static-debian12:latest
COPY --from=builder /usr/local/bin/function /usr/local/bin/function
ENTRYPOINT ["function"]
