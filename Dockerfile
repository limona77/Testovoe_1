FROM golang:1.22.0

WORKDIR /app


COPY . .

RUN go build -o task.exe main.go

RUN go mod tidy


CMD ["./task.exe", "test_file.txt"]