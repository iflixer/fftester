# build environment
FROM golang:1.21.1 as build-env
WORKDIR /server
COPY src/go.mod ./
RUN go mod download
COPY src src
WORKDIR /server/src
RUN CGO_ENABLED=0 GOOS=linux go build -o /server/build/httpserver .

FROM linuxserver/ffmpeg
WORKDIR /app

COPY --from=build-env /server/build/httpserver /app/videoconverter

#ENV GITHUB-SHA=<GITHUB-SHA>

ENTRYPOINT [ "/app/videoconverter" ]
#ENTRYPOINT [ "ls", "-la", "/app/httpserver" ]
