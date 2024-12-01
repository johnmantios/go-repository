FROM postgres:alpine
COPY ./../migrations/*.up.sql /docker-entrypoint-initdb.d/
EXPOSE 5432