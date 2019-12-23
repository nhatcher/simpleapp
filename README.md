Simple Login App
================

Description
-----------

This is meant as a simple exercise for Carlitos to learn about:

* Programming
* Go
* SQL (SQLite and PostgreSQL)
* Javascript
* HTML/CSS
* UI/UX design

The idea is to build a simple website that is is password protected.
To visit the website you need to have a username and a password.
There is a special type of user (admin) that can add/remove/edit users from a particular endpoint

Build Instructions
------------------

To build the app and run the app:

    $ make
    $ ./simpleapp

Then visit http://localhost:1312/

To remove all artifacts (in particular the database data) run:

$ make clean

Dependencies
------------

Aside from the Go programming language you will need the two libraries:

$ go get github.com/mattn/go-sqlite3
$ go get golang.org/x/crypto/bcrypt

Note that there are no dependencies to run the server once it is built.
