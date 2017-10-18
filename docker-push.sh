#/bin/bash

NAME=${1:-${PWD##*/}} #name from argument or current directory name
NAME=$(echo $NAME | tr "[:upper:]" "[:lower:]")

VERSION=${2:-'latest'}
VERSION=$(echo $VERSION | tr "[:upper:]" "[:lower:]")

REGISTRY=registry.acamica.com

for PUSH_VERSION in $VERSION `git rev-parse --short HEAD` `jq -r '.version' package.json` `git rev-parse --abbrev-ref HEAD`; do
    echo "push $NAME rev $PUSH_VERSION"
    docker tag $NAME:$VERSION $REGISTRY/$NAME:$PUSH_VERSION
    docker push $REGISTRY/$NAME:$PUSH_VERSION
done

