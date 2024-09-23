#!/bin/sh

set -e  # will exit the script immediately if any of the command fails 

echo "Starting the application..."

echo "Sourcing app.env" 
source app.env 
echo "$DRIVER_SOURCE"

echo "Running Migration"
migrate -path=db/migration -database=$DRIVER_SOURCE up

exec "$@" # hand over the execution to command passed, In our case, ./app/simple_bank app 