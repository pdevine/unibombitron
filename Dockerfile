FROM golang:1.17-alpine3.14
WORKDIR /project
COPY *.go ./
COPY go.* ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o bombitron *.go

FROM scratch
COPY --from=0 /project/bombitron /bombitron
COPY *.png ./
ENV COLORTERM="truecolor"
ENV TERM="xterm-256color"
ENTRYPOINT ["/bombitron"]
