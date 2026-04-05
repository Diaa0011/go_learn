# Make sure to specify the same Go version as the one in the go.mod file.
FROM golang:1.26.1

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o svr ./

RUN chmod +x /app/svr

EXPOSE 8080
CMD [ "./svr" ]