CREATE ROLE api 
    WITH CREATEDB 
    LOGIN 
    PASSWORD 'open-db';

CREATE DATABASE "api-db" 
    OWNER api 
    ENCODING 'UTF8';
