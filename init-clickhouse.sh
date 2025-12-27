#!/bin/sh
set -euo pipefail

CH_HOST="${CH_HOST:-clickhouse}"
CH_PORT="${CH_PORT:-9000}"   # clickhouse-client ходит на native порт
CH_DB_CHECK="${CH_DB_CHECK:-streaming}"
INIT_DIR="${INIT_DIR:-/init}"

echo "[init-ch] host=${CH_HOST} port=${CH_PORT} init_dir=${INIT_DIR}"

echo "[init-ch] waiting for clickhouse..."
until clickhouse-client --host "${CH_HOST}" --port "${CH_PORT}" --query "SELECT 1" >/dev/null 2>&1; do
  sleep 1
done
echo "[init-ch] clickhouse is ready"

echo "[init-ch] listing init dir:"
ls -la "${INIT_DIR}"

# Собираем файлы *.sql
set +e
FILES=$(ls -1 "${INIT_DIR}"/*.sql 2>/dev/null)
set -e

if [ -z "${FILES:-}" ]; then
  echo "[init-ch] ERROR: no *.sql files found in ${INIT_DIR}"
  exit 1
fi

echo "[init-ch] applying sql files in order:"
echo "${FILES}"

for f in ${FILES}; do
  echo "[init-ch] apply: ${f}"
  clickhouse-client --host "${CH_HOST}" --port "${CH_PORT}" --multiquery < "${f}"
done

echo "[init-ch] verifying database exists: ${CH_DB_CHECK}"
clickhouse-client --host "${CH_HOST}" --port "${CH_PORT}" --query "EXISTS DATABASE ${CH_DB_CHECK}"

echo "[init-ch] OK"
