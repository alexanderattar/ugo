#!/bin/sh

# Start the Docker container
docker run -d --rm --name ujo-postgres \
    -p 5432:5432 \
    -e POSTGRES_USER=docker \
    -e POSTGRES_PASSWORD=docker \
    -e POSTGRES_DB=ujo \
    postgres
# Uncomment these for persistent volumes
#    -v /tmp/data:/var/lib/postgresql/data \
#    -v /tmp/pgrun:/var/run/postgresql \

# @@TODO: this is hacky
sleep 5

# Apply migrations and load fixtures
#DATABASE_USERNAME=docker DATABASE_PASSWORD=docker go run db/loaddata.go


