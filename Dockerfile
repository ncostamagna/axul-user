FROM golang:1.22

WORKDIR $GOPATH/bin

COPY main .

RUN mkdir log 

EXPOSE 8082

CMD ["main"]