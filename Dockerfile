# Base image for building the go project
FROM golang:1.14-alpine AS build

# Updates the repository and installs git
RUN apk update && apk upgrade && \
    apk add --no-cache git

# Switches to /tmp/app as the working directory, similar to 'cd'
WORKDIR /tmp/app

COPY . .

# Builds the current project to a binary file called api
# The location of the binary file is /tmp/app/out/api
RUN GOOS=linux go build -o ./out/api .

#########################################################

# The project has been successfully built and we will use a
# lightweight alpine image to run the server 
FROM alpine:latest

# Adds CA Certificates to the image
RUN apk add ca-certificates

# Copies the binary file from the BUILD container to /app folder
COPY --from=build /tmp/app/out/api /app/api

# Switches working directory to /app
WORKDIR "/app"

# Exposes the 5000 port from the container
EXPOSE 5000

# Runs the binary once the container starts
CMD ["./api"]


# FROM golang:1.15

#RUN mkdir -p $GOPATH/src/github.com/ncostamagna/axul_user
#WORKDIR $GOPATH/src/github.com/ncostamagna/axul_user
#COPY . .
#RUN ls

#ARG DATABASE_HOST
#ARG DATABASE_USER
#ARG DATABASE_PASSWORD 
#ARG DATABASE_NAME
#ARG DATABASE_PORT
#ARG DATABASE_DEBUG
#ARG DATABASE_MIGRATE
#ARG APP_PORT
#ARG APP_URL

#ENV DATABASE_HOST $DATABASE_HOST
#ENV DATABASE_USER $DATABASE_USER
#ENV DATABASE_PASSWORD $DATABASE_PASSWORD
#ENV DATABASE_NAME $DATABASE_NAME
#ENV DATABASE_PORT $DATABASE_PORT
#ENV DATABASE_DEBUG $DATABASE_DEBUG
#ENV DATABASE_MIGRATE $DATABASE_MIGRATE
#ENV APP_PORT $APP_PORT
#ENV APP_URL $APP_URL

#RUN go get -d -v ./... 
#RUN go install -v ./...
#EXPOSE 8082
#EXPOSE 50055

#CMD ["axul_user"]