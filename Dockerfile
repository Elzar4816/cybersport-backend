FROM golang:1.23-alpine

WORKDIR /app

COPY . .
RUN rm -rf frontend

RUN go mod tidy && go build -o server .

EXPOSE 8000
CMD ["./server"]
