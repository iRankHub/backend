FROM postgres:13-alpine
COPY scripts/init-user-db.sh /docker-entrypoint-initdb.d/
COPY internal/database/postgres/migrations /docker-entrypoint-initdb.d/