#!/bin/zsh

CONTAINER_NAME="game-server-example-mysql"
MYSQL_USER="root"
CHANGE_GENERAL_LOG_FILE_QUERY="set global general_log_file='/var/lib/mysql/localhost.log';"
ENABLE_GENERAL_LOG_QUERY="set global general_log=1;"

docker exec -it ${CONTAINER_NAME} touch /var/lib/mysql/localhost.log
docker exec -it ${CONTAINER_NAME} chmod 0666 /var/lib/mysql/localhost.log
docker exec -it ${CONTAINER_NAME} mysql -u ${MYSQL_USER} ${MYSQL_DB} -e "${CHANGE_GENERAL_LOG_FILE_QUERY}"
docker exec -it ${CONTAINER_NAME} mysql -u ${MYSQL_USER} ${MYSQL_DB} -e "${ENABLE_GENERAL_LOG_QUERY}"
docker exec -it ${CONTAINER_NAME} truncate -s 0 /var/lib/mysql/localhost.log
