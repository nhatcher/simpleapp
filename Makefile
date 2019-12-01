all: database.sqlite
	go build

database.sqlite:
	sqlite3 database.sqlite < database.sql
	sqlite3 database.sqlite < data.sql

clean:
	rm -rf database.sqlite
	rm -rf simpleapp
