FROM ubuntu:14.04

# Setup MongoDB sources
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv 7F0CEB10
RUN echo "deb http://repo.mongodb.org/apt/ubuntu trusty/mongodb-org/3.0 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-3.0.list

RUN apt-get update && apt-get -y upgrade
RUN apt-get install -y golang git mongodb-org-tools curl

# Set GOPATH
ENV GOPATH /go

# Grab Source
COPY . /go/src/github.com/scheedule/coursestore
COPY start.sh .

# Grab project dependencies
RUN cd /go/src/github.com/scheedule/coursestore/scrape && go get ./...
RUN cd /go/src/github.com/scheedule/coursestore/serve && go get ./...

# Build Project
RUN cd /go/src/github.com/scheedule/coursestore/scrape && go install
RUN cd /go/src/github.com/scheedule/coursestore/serve && go install

ENTRYPOINT ["bash", "start.sh"]
