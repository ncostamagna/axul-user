FROM golang:1.15

RUN mkdir -p $GOPATH/src/github.com/ncostamagna/axul_user
WORKDIR $GOPATH/src/github.com/ncostamagna/axul_user
COPY . .
RUN ls

ARG DATABASE_PASSWORD 
ARG DATABASE_NAME

ENV DATABASE_HOST 'la pija negra'
ENV DATABASE_USER $DATABASE_USER
ENV DATABASE_PASSWORD $DATABASE_PASSWORD
ENV DATABASE_NAME DATABASE_NAME

RUN go get -d -v ./... 
RUN go install -v ./...
EXPOSE 8082
EXPOSE 50055

CMD ["axul_user"]