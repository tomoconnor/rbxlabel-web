FROM golang:1.18 as builder

WORKDIR /go/src/github.com/tomoconnor/rbxlabel-web

COPY go.mod go.sum main.go api.go ./

RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rbxlabel-web .

FROM scratch
WORKDIR /
COPY --from=builder /go/src/github.com/tomoconnor/rbxlabel-web/rbxlabel-web .
ENV PORT=7800
EXPOSE 7800
ENTRYPOINT ["./rbxlabel-web"]
