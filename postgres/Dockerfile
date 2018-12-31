FROM postgres:11.1
COPY init_kaiser.sql /docker-entrypoint-initdb.d/1.sql
COPY init_kong.sh /docker-entrypoint-initdb.d/2.sh
CMD ["postgres"]