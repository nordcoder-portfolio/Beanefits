#!/bin/sh
set -e

BROKERS="${KAFKA_BROKERS:-redpanda:19092}"
TOPIC="${KAFKA_TOPIC:-ints}"
PARTITIONS="${KAFKA_PARTITIONS:-1}"
REPLICAS="${KAFKA_REPLICAS:-1}"

echo "[init] Waiting for Redpanda at ${BROKERS}..."
until rpk cluster info -X brokers="${BROKERS}" >/dev/null 2>&1; do
  sleep 1
done

echo "[init] Redpanda is up. Ensuring topic '${TOPIC}' exists..."

if rpk topic describe "${TOPIC}" -X brokers="${BROKERS}" >/dev/null 2>&1; then
  echo "[init] Topic '${TOPIC}' already exists."
else
  rpk topic create "${TOPIC}" \
    --partitions "${PARTITIONS}" \
    --replicas "${REPLICAS}" \
    -X brokers="${BROKERS}"
  echo "[init] Topic '${TOPIC}' created."
fi

echo "[init] Current topics:"
rpk topic list -X brokers="${BROKERS}"
