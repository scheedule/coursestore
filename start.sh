#!/bin/sh
if [ "$1" == "scrape" ]
then
  /go/bin/scrape mongo
else
  curl --insecure https://www.scheedule.com/dump.tar.gz | tar xz --
  mongorestore --host mongo --collection classes
fi

/go/bin/serve
