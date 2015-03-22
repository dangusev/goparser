FROM golang:1.4.2
MAINTAINER Daniil Gusev <dangusev92@gmail.com>
EXPOSE 8080

RUN mkdir -p /go/src/github.com/dangusev/goparser
ADD . /go/src/github.com/dangusev/goparser
WORKDIR /go/src/github.com/dangusev/goparser

# Install apt dependencies
RUN apt-get update
RUN apt-get install -y apt-utils
RUN apt-get install -y pkg-config
RUN apt-get install -y libxml2 libxml2-dev git npm
RUN ln -s /usr/bin/nodejs /usr/bin/node

# Install Bower and get static files
RUN npm install -g bower
RUN bower install --allow-root
RUN cd static
RUN ln -s ../bower_components/angular/angular.min.js .
RUN ln -s ../bower_components/angular-ui-router/release/angular-ui-router.min.js .

# Build and install Go application
RUN go-wrapper download && go-wrapper install

ENTRYPOINT /go/bin/goparser