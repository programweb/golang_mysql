This page looks good w/o docker-compose.yml        https://www.configserverfirewall.com/docker/run-mysql-docker-containers/
using:                                             https://stackoverflow.com/questions/48808069/how-to-create-database-in-database-docker-container
- - -
$ docker-compose build
database uses an image, skipping

$ docker-compose up -d database
$ docker ps
$ docker exec -it mysql bash
# mysql -u programweb -pABC                (or just mysql -u root    but it should work with pw now)
mysql> SHOW databases;
mysql> USE healthdata;
mysql> SHOW tables;             -- it shows artists
mysql> SELECT * FROM cause;

mysql -uprogramweb -hmysql -P33061 -p

