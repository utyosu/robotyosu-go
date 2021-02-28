fmt:
	go fmt ./...

tidy:
	go mod tidy

build-local: fmt tidy
	go build -o bin/robotyosu-local -tags="local"

build-production:
	GOOS=linux GOARCH=amd64 go build -o bin/robotyosu-production -tags="production"

run-local: build-local
	sleep 0.5
	./bin/robotyosu-local

deploy-production: build-production
	scp bin/robotyosu-production production:/tmp
	ssh production " \
		mv /tmp/robotyosu-production /home/ec2-user && \
		supervisorctl restart robotyosu \
	"

start-production:
	ssh production "supervisorctl start robotyosu"

stop-production:
	ssh production "supervisorctl stop robotyosu"

reset-db-local:
	sudo mysql -u ${DATABASE_USER_LOCAL} -e " \
		DROP DATABASE IF EXISTS ${DATABASE_NAME_LOCAL}; \
		CREATE DATABASE ${DATABASE_NAME_LOCAL}; \
	"
	@make migrate-db-up-local
	if [ -e "test.sql" ]; then sudo mysql -u ${DATABASE_USER_LOCAL} -h ${DATABASE_HOST_LOCAL} -P ${DATABASE_PORT_LOCAL} ${DATABASE_NAME_LOCAL} < test.sql; fi

migrate-db-up-local:
	migrate -path db/migrations -database "mysql://${DATABASE_USER_LOCAL}:@tcp(${DATABASE_HOST_LOCAL}:${DATABASE_PORT_LOCAL})/${DATABASE_NAME_LOCAL}" up
	mysqldump -u ${DATABASE_USER_LOCAL} -h ${DATABASE_HOST_LOCAL} -P ${DATABASE_PORT_LOCAL} ${DATABASE_NAME_LOCAL} -d --skip-comments > db/schema.sql

migrate-db-down-local:
	migrate -path db/migrations -database "mysql://${DATABASE_USER_LOCAL}:@tcp(${DATABASE_HOST_LOCAL}:${DATABASE_PORT_LOCAL})/${DATABASE_NAME_LOCAL}" down 1
	mysqldump -u ${DATABASE_USER_LOCAL} -h ${DATABASE_HOST_LOCAL} -P ${DATABASE_PORT_LOCAL} ${DATABASE_NAME_LOCAL} -d --skip-comments > db/schema.sql

migrate-db-production:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	migrate -path db/migrations -database "mysql://${DATABASE_USER_PRODUCTION}:${DATABASE_PASSWORD_PRODUCTION}@tcp(${DATABASE_HOST_PRODUCTION}:${DATABASE_PORT_PRODUCTION})/${DATABASE_NAME_PRODUCTION}" up

provisioning-production:
	scp -r supervisor/ production:/tmp/
	ssh production " \
		sudo yum install python3 -y && \
		sudo pip3 install supervisor && \
		sudo rm -rf /etc/supervisor && \
		sudo mv /tmp/supervisor /etc && \
		killall supervisord | true && \
		supervisord \
	"
