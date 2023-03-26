FROM golang:1.19

WORKDIR $GOPATH/bin

COPY main .

EXPOSE 8082

CMD ["main"]