TIMESTAMP := $(shell TZ='Asia/Tokyo' date '+%Y%m%d%H%M')
BACKUP_DIR := db/backups
BACKUP_FILE_PREFIX := dump_rbc_production_
LD_FLAGS := \
  -X 'github.com/utyosu/robotyosu-go/app.commitHash=$(shell git log --pretty=format:%H -n 1)' \
  -X 'github.com/utyosu/robotyosu-go/app.buildDatetime=$(shell TZ='Asia/Tokyo' date '+%Y-%m-%d %H:%M:%S JST')'

check-go-version:
	./check-go-version.sh

fmt:
	go fmt ./...

tidy:
	go mod tidy

build-local: fmt tidy check-go-version
	go build -ldflags="${LD_FLAGS}" -o bin/robotyosu-local -tags="local"

build-production: fmt tidy check-go-version
	GOOS=linux GOARCH=amd64 go build -ldflags="${LD_FLAGS}" -o bin/robotyosu-production -tags="production"

run-local: build-local
	sleep 0.5
	./bin/robotyosu-local

deploy-production: build-production
	scp bin/robotyosu-production production:/tmp
	ssh production " \
		mv /tmp/robotyosu-production /home/ec2-user && \
		sudo /usr/local/bin/supervisorctl restart robotyosu \
	"

start-production:
	ssh production "supervisorctl start robotyosu"

stop-production:
	ssh production "supervisorctl stop robotyosu"

reset-db-local:
	sudo mysql -u ${RBC_DATABASE_USER_LOCAL} -e " \
		DROP DATABASE IF EXISTS ${RBC_DATABASE_NAME_LOCAL}; \
		CREATE DATABASE ${RBC_DATABASE_NAME_LOCAL}; \
	"
	@make migrate-db-up-local
	if [ -e "test.sql" ]; then sudo mysql -u ${RBC_DATABASE_USER_LOCAL} -h ${RBC_DATABASE_HOST_LOCAL} -P ${RBC_DATABASE_PORT_LOCAL} ${RBC_DATABASE_NAME_LOCAL} < test.sql; fi

migrate-db-up-local:
	migrate -path db/migrations -database "mysql://${RBC_DATABASE_USER_LOCAL}:@tcp(${RBC_DATABASE_HOST_LOCAL}:${RBC_DATABASE_PORT_LOCAL})/${RBC_DATABASE_NAME_LOCAL}" up
	mysqldump -u ${RBC_DATABASE_USER_LOCAL} -h ${RBC_DATABASE_HOST_LOCAL} -P ${RBC_DATABASE_PORT_LOCAL} ${RBC_DATABASE_NAME_LOCAL} -d --skip-comments --no-tablespaces | sed 's/ AUTO_INCREMENT=[0-9]*//g' > db/schema.sql

migrate-db-down-local:
	migrate -path db/migrations -database "mysql://${RBC_DATABASE_USER_LOCAL}:@tcp(${RBC_DATABASE_HOST_LOCAL}:${RBC_DATABASE_PORT_LOCAL})/${RBC_DATABASE_NAME_LOCAL}" down 1
	mysqldump -u ${RBC_DATABASE_USER_LOCAL} -h ${RBC_DATABASE_HOST_LOCAL} -P ${RBC_DATABASE_PORT_LOCAL} ${RBC_DATABASE_NAME_LOCAL} -d --skip-comments --no-tablespaces | sed 's/ AUTO_INCREMENT=[0-9]*//g' > db/schema.sql
