-- User: api_user
-- DROP USER api_user;

CREATE USER api_user WITH
    LOGIN
    NOSUPERUSER
    INHERIT
    CREATEDB
    CREATEROLE
    NOREPLICATION;