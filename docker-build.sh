#/bin/bash

NAME=${1:-${PWD##*/}} #name from argument or current directory name
NAME=$(echo $NAME | tr "[:upper:]" "[:lower:]")

VERSION=${2:-'latest'}
VERSION=$(echo $VERSION | tr "[:upper:]" "[:lower:]")

docker build -t $NAME:$VERSION .

