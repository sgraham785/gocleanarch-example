-- DROP DATABASE gcarch_example;
-- CREATE DATABASE gcarch_example;

-- use gcarch_example;

CREATE TABLE IF NOT EXISTS "user" (
  id varchar(50),
  email varchar(255),
  password varchar(255),
  first_name varchar(100),
  last_name varchar(100),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id));

CREATE TABLE IF NOT EXISTS book (
  id varchar(50),
  title varchar(255),
  author varchar(255),
  pages integer,
  quantity integer,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (id));

CREATE TABLE IF NOT EXISTS book_user (
  user_id varchar(50),
  book_id varchar(50),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  PRIMARY KEY (user_id, book_id));

