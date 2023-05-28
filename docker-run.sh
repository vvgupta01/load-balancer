#!/bin/bash
image="go_load_balancer"
container="load_balancer"


if [ $(docker ps -a -f "name=$container" -q) ]
then
    echo "Removing $container container..."
    docker rm -f $container
fi

if [ $(docker images $image -q) ]
then
    echo "Removing $image image..."
    docker image rm -f $image
fi

docker build -t $image .
docker run --network host -it --name $container $image