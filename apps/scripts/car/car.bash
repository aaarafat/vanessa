#!/bin/bash

usage() { echo "Usage: $0 [-p|--port <port>]  [-i|--id <id>]" 1>&2; exit 1; }

if [ $# -eq 0 ]; then
    usage
fi
while [ "$1"  != "" ]; do
  FLAG=`echo $1 | awk -F= '{print $1}'`
  shift
  VALUE=`echo $1 | awk -F= '{print $1}'`

  case $FLAG in
    -p|--port)
      PORT=$VALUE
      ;;
    -i|--id)
      ID=$VALUE
      ;;
    *)
      usage
      exit 1
  esac
  shift
done

# if car doesn't exist in current directory
# go to directory scripts/car
if [ ! -d "car.go" ]; then
  cd ./apps/scripts/car
fi
echo $PORT
echo $ID
socat TCP4-LISTEN:$PORT,fork,reuseaddr UNIX-CONNECT:/tmp/car$ID.ui.socket &
# if go is not in path add it
if [ ! -x "$(command -v go)" ]; then
  export PATH=$PATH:/usr/local/go/bin
fi
go build -o car

./car -id $ID -s /tmp/car$ID.ui.socket



