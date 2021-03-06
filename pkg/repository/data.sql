DROP DATABASE IF EXISTS reservation;
CREATE DATABASE reservation;

USE reservation;

DROP TABLE IF EXISTS users;
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(32) NOT NULL UNIQUE,
    password VARCHAR(64) NOT NULL,
    status INTEGER DEFAULT 0
);

INSERT INTO users (
    username,
    password) 
VALUES (
    "admin",
    "admin"
);

DROP TABLE IF EXISTS halls;
CREATE TABLE halls (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(64),
    capacity INTEGER,
    reserved INTEGER DEFAULT 0
);

DROP TABLE IF EXISTS reserves;
CREATE TABLE reserves (
    id INTEGER PRIMARY KEY AUTO_INCREMENT,
    user_id INTEGER,
    hall_name VARCHAR(64)
);
