FROM golang:latest 
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
#RUN set -x && \ 
#    go get -d -v . && \
#    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
COPY . .
RUN go build -o app .
#FROM scratch
#WORKDIR /root/
#COPY --from=0 /go/src/github.com/user/app .
EXPOSE 3030
CMD ["./app"]  
