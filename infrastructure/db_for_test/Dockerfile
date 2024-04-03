FROM mysql:8
LABEL authors="asragi"

ENV MYSQL_ROOT_PASSWORD ringo
EXPOSE 3306
ADD ./my.cnf /etc/mysql/conf.d/my.cnf
ADD ./init /docker-entrypoint-initdb.d
CMD ["mysqld"]
