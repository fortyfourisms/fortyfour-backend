#!/bin/sh
set -e

# Wait for MySQL
echo "Waiting for MySQL..."
while ! nc -z $DB_HOST $DB_PORT; do
  sleep 1
done
echo "MySQL is up!"

# Wait for Redis
echo "Waiting for Redis..."
while ! nc -z $REDIS_HOST $REDIS_PORT; do
  sleep 1
done
echo "Redis is up!"

# Run migrations (optional)
# echo "Running migrations..."
# ./migrate -path ./migrations -database "mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}" up

# Start application
echo "Starting application..."
exec "$@"