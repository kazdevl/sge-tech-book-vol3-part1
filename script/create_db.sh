#!/bin/zsh

CONTAINER_NAME="game-server-example-mysql"
MYSQL_USER="root"
CREATE_DB_QUERY="create database game_server_example;"

docker exec -it ${CONTAINER_NAME} mysql -u ${MYSQL_USER} -e "${CREATE_DB_QUERY}"
