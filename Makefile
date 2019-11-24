all: database.sqlite
	go build

database.sqlite:
	sqlite3 database.sqlite < database.sql

clean:
	rm database.sqlite
	rm webapp
