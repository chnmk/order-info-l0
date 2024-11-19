FROM golang:1.22-alpine3.19

WORKDIR /orders

COPY . ./

RUN go mod tidy 
