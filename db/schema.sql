DROP DATABASE IF EXISTS tasker CASCADE;

CREATE DATABASE IF NOT EXISTS tasker;

SET DATABASE = tasker;

CREATE TABLE IF NOT EXISTS tasks (
     taskId SERIAL NOT NULL PRIMARY KEY,
     skills VARCHAR NOT NULL UNIQUE,
     priority VARCHAR NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS agents (
     agentId SERIAL NOT NULL PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS skills (
    skillId SERIAL NOT NULL PRIMARY KEY,
    name VARCHAR NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS skillmapping (
      skillMapId SERIAL NOT NULL PRIMARY KEY,
      skill VARCHAR NOT NULL,
      agentId INTEGER NOT NULL
);

INSERT INTO skillmapping(skillMapId, skill, agentId) VALUES
(1, 'skill1', 1),
(2, 'skill1', 2),
(3, 'skill1', 3),
(4, 'skill2', 1),
(5, 'skill2', 3),
(6, 'skill3', 1),
(7, 'skill3', 2);


