-- sqlite3 database.sqlite < data.sql

INSERT INTO USERS (
		first_name,
		last_name,
		email,
		username,
		password,
		usertype_id
	) VALUES (
    'John',
		'Smith',
		'jsmith@example.com',
		'jsmith',
		'$2a$10$iMbDo5snmxqm.SBpN7JJD.CABTbpimrUps7XfAe5jQMN.m4c0uq6u',
		(SELECT usertype_id FROM usertype WHERE type="user")
  );
  
  INSERT INTO USERS (
		first_name,
		last_name,
		email,
		username,
		password,
		usertype_id
	) VALUES (
    'Juan',
		'Pablo',
		'jpablo@example.com',
		'jpablo',
		'$2a$10$iMbDo5snmxqm.SBpN7JJD.CABTbpimrUps7XfAe5jQMN.m4c0uq6u',
		(SELECT usertype_id FROM usertype WHERE type="root")
  );