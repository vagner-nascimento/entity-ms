#########
# build #
#########
FROM golang:1.16 as build
WORKDIR /app
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO_ENV=prod

COPY ["./go.mod", "./go.sum", "./main.go", "./"]
COPY ./src ./src

RUN go mod download
RUN go build -o main .

##########
# deploy #
##########
FROM gcr.io/distroless/base-debian11
WORKDIR /app
COPY --from=build /app .
ENTRYPOINT [ "/app/main" ]
