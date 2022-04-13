FROM golang:latest AS builder

WORKDIR /build

COPY go.mod go.sum ./
    
RUN go mod download
    
COPY . .
    
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Starting a new stage from scratch 
    
FROM alpine:latest
    
RUN apk --no-cache add ca-certificates
    
#WORKDIR /root/
    
#COPY --from=builder /app/main .

WORKDIR /build

COPY --from=builder /build/main /build/main
    
EXPOSE 10000 10000
    
CMD ["./main"]

