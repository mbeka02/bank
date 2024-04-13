#!/bin/sh

set -e

echo "run db migration"
/app/goose -dir /app/schema postgres "$DB_URL" up
echo "start app"
exec "$@"

