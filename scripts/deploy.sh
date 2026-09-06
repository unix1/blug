#!/bin/sh
set -e
cd "$(dirname "$0")/.."

set -a
. ./.env
set +a

if [ -z "$S3_PROVIDER" ] || [ -z "$S3_REGION" ] || [ -z "$S3_ENDPOINT" ] || [ -z "$S3_ACCESS_KEY_ID" ] || [ -z "$S3_SECRET_ACCESS_KEY" ] || [ -z "$S3_BUCKET" ]; then
  echo "S3_PROVIDER, S3_REGION, S3_ENDPOINT, S3_ACCESS_KEY_ID, S3_SECRET_ACCESS_KEY, and S3_BUCKET must be set in .env" >&2
  exit 1
fi

rclone sync public/ ":s3:${S3_BUCKET}" \
  --config /dev/null \
  --s3-provider "$S3_PROVIDER" \
  --s3-region "$S3_REGION" \
  --s3-endpoint "$S3_ENDPOINT" \
  --s3-access-key-id "$S3_ACCESS_KEY_ID" \
  --s3-secret-access-key "$S3_SECRET_ACCESS_KEY" \
  -v
