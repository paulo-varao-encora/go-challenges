-- create database go_challenges;
-- source /home/pvarao/Documents/go/challenges/create-tables.sql

DROP TABLE IF EXISTS tasks;
CREATE TABLE tasks (
  ID        INT AUTO_INCREMENT NOT NULL,
  Name      VARCHAR(128) NOT NULL,
  Completed TINYINT NOT NULL,
  PRIMARY KEY (`ID`)
);

INSERT INTO tasks
  (Name, Completed)
VALUES
  ('Pay bills', 1),
  ('Walk the dog', 0),
  ('Buy groceries', 0),
  ('Exercise', 1);