FROM mysql:latest

COPY ./deploy ./docker-entrypoint-initdb.d/

WORKDIR /docker-entrypoint-initdb.d/

ENV MYSQL_DATABASE=notedb
ENV MYSQL_ROOT_PASSWORD=mysql

EXPOSE 3306

# docker build -t notedb .
# docker run -p 3306:3306 --rm notedb:latest 