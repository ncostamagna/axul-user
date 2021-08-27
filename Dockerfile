FROM golang:1.15

RUN mkdir -p $GOPATH/src/github.com/ncostamagna/axul_user
WORKDIR $GOPATH/src/github.com/ncostamagna/axul_user
COPY . .
RUN ls

RUN go get -d -v ./... 
RUN go install -v ./...
EXPOSE 8082

CMD ["axul_user"]