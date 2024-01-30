FROM mysql:latest

COPY ./deploy ./docker-entrypoint-initdb.d/

WORKDIR /docker-entrypoint-initdb.d/

ENV MYSQL_DATABASE=notedb
ENV MYSQL_ROOT_PASSWORD=mysql

EXPOSE 3306