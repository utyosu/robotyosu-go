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
	sudo mysql -u root -e " \
		DROP DATABASE IF EXISTS robotyosu_local; \
		CREATE DATABASE robotyosu_local; \
	"
	sudo mysql -u root robotyosu_local < schema/schema.sql
	if [ -e "test.sql" ]; then sudo mysql -u root robotyosu_local < test.sql; fi

reset-db-production:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]
	mysql -h ${PRODUCTION_HOST} -u ${PRODUCTION_USER} -p -e " \
		DROP DATABASE IF EXISTS robotyosu_production; \
		CREATE DATABASE robotyosu_production; \
	"
	mysql -h ${PRODUCTION_HOST} -u ${PRODUCTION_USER} robotyosu_production -p < schema/schema.sql

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
