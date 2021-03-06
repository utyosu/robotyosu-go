set -eu

TEST_DB_NAME="robotyosu_test"
TEST_DB_USER="robotoysu_test"
TEST_DB_HOST="127.0.0.1"
TEST_DB_PORT="3306"

function up() {
  # 一時的にエラーを許容する
  set +e
  result=`migrate -path db/migrations -database "mysql://${TEST_DB_USER}:@tcp(${TEST_DB_HOST}:${TEST_DB_PORT})/${TEST_DB_NAME}" up ${1} 2>&1`
  return_code=$?
  set -e
  # エラー許容ここまで

  # マイグレーションファイルが存在しない（=最新バージョン）のときは成功で終了する
  if [ "${result}" = "error: file does not exist" ]; then
    echo "Success migration test"
    exit 0
  fi

  echo "${result}"

  # 終了コードが0でなければエラーにする
  if [ "${return_code}" != "0" ]; then
    exit ${return_code}
  fi
}

function down() {
  migrate -path db/migrations -database "mysql://${TEST_DB_USER}:@tcp(${TEST_DB_HOST}:${TEST_DB_PORT})/${TEST_DB_NAME}" down ${1}
}

function dump() {
  mysqldump -u ${TEST_DB_USER} -h ${TEST_DB_HOST} -P ${TEST_DB_PORT} ${TEST_DB_NAME} --skip-comments --no-tablespaces > db/migrations/tmp/${1}.sql
}

function import() {
  path="db/migrations/test_data/${1}.sql"
  if [ -f ${path} ]; then
    echo "Import test data: ${path}"
    mysql -u ${TEST_DB_USER} -h ${TEST_DB_HOST} -P ${TEST_DB_PORT} ${TEST_DB_NAME} < ${path}
  fi
}

function checkDiff() {
  result_diff=`diff -u db/migrations/tmp/${1}_up.sql db/migrations/tmp/${1}_down.sql | cat`
  if [ -n "${result_diff}" ]; then
    echo "${result_diff}"
    exit 1
  fi
}

sudo mysql -e " \
  DROP DATABASE IF EXISTS ${TEST_DB_NAME}; \
  CREATE DATABASE ${TEST_DB_NAME}; \
  CREATE USER IF NOT EXISTS ${TEST_DB_USER};
  GRANT ALL ON ${TEST_DB_NAME}.* TO ${TEST_DB_USER};
"

rm -rf db/migrations/tmp
mkdir -p db/migrations/tmp
version=0

# schema_migrationsテーブルを作るために1回上げ下げしておく
up 1
down 1

while true
do
  import ${version}
  dump "${version}_up"
  up 1
  down 1
  dump "${version}_down"
  checkDiff ${version}
  version=$((version+1))
  up 1
done
