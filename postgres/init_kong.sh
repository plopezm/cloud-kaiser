#!/usr/bin/env bash
export PGUSER=postgres
psql <<- EOSQL
    CREATE DATABASE kongdb;
    GRANT ALL PRIVILEGES ON DATABASE kongdb TO postgres;
EOSQL