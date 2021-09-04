FROM golang:alpine3.12
WORKDIR /project
COPY bombitron.go .
COPY go.* ./
RUN apk add --no-cache git
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o bombitron bombitron.go

FROM scratch
COPY --from=0 /project/bombitron /bombitron
COPY *.png ./
ENV COLORTERM="truecolor"
ENV TERM="xterm-256color"
ENTRYPOINT ["/bombitron"]
