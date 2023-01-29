FROM golang:1.19

RUN mkdir -p $GOPATH/src/github.com/ncostamagna/axul_user
WORKDIR $GOPATH/src/github.com/ncostamagna/axul_user

COPY . .
RUN ls

RUN go install
EXPOSE 8082

CMD ["axul_user"]