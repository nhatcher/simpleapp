-- sqlite3 database.sqlite < database.sql

CREATE TABLE USERS (
  user_id INTEGER PRIMARY KEY,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  usertype_id INTEGER,
  FOREIGN KEY (usertype_id) REFERENCES USERTYPE (usertype_id)
);

CREATE TABLE SESSIONS (
  session_id INTEGER PRIMARY KEY,
  user_id INTEGER,
  session_hash TEXT NOT NULL,
  create_date DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES USERS (user_id)
);

CREATE TABLE USERTYPE (
  usertype_id INTEGER PRIMARY KEY,
  type TEXT NOT NULL,
  description TEXT NOT NULL
);

INSERT INTO USERTYPE (
  type,
  description
) VALUES (
  "user",
  "Normal user"
);

INSERT INTO USERTYPE (
  type,
  description
) VALUES (
  "root",
  "Admin user"
);
