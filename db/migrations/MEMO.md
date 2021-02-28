# Migrations

## Install migrate

curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
mv ./migrate.linux-amd64 /usr/bin/migrate

## Create

migrate create -ext sql -dir db/migrations -seq name
