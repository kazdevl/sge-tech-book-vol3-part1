.PHONY: api_gen
api_gen:
	oapi-codegen -package schema openapi.yaml > pkg/handler/schema/api.gen.go

.PHONY: model_gen
model_gen:
	sqlboiler mysql --wipe --pkgname datamodel --output pkg/infra/mysql/datamodel

.PHONY: docker_up_db
docker_up_db:
	docker-compose -f docker-compose.yml up --build -d mysql

.PHONY: exec_db_setting
exec_db_setting:
	sh ./script/create_db.sh
	sh ./script/general_log_setting.sh
	make migrate_local

.PHONY: exec_general_log_setting
exec_general_log_setting:
	sh ./script/general_log_setting.sh

.PHONY: docker_up
EnableProposalMethod=1
docker_up:
	docker-compose -f docker-compose.yml --compatibility up --build -d

.PHONY: migrate_local
DSN:=mysql://root:@tcp(127.0.0.1:13502)/game_server_example
migrate_local:
	yes | migrate -path ./ddl/migration/ -database "$(DSN)" drop
	yes | migrate -path ./ddl/migration/ -database "$(DSN)" up

.PHONY: save_result
FILEPATH=./result.log
save_result:
	sh ./script/save_result.sh $(FILEPATH)

.PHONY: docker_down
docker_down:
	docker-compose -f docker-compose.yml down
