#!/bin/zsh

CONTAINER_NAME="game-server-example-mysql"

result=`docker exec -it ${CONTAINER_NAME} cat /var/lib/mysql/localhost.log`
echo "${result}" > $1
docker exec -it ${CONTAINER_NAME} truncate -s 0 /var/lib/mysql/localhost.log
