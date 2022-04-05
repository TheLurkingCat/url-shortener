# URL Shortener

## Description

Use shorten ID as sqlite database's primary key. 
When accessing a shorten URL, it will query the database.

If no row is found, it replies a 404 Not Found message. 
Otherwise it will response with a 302 Moved Temporarily which redirect to original URL.

## Dependencies

- [gin](github.com/gin-gonic/gin) \- A fast an easy to use web framework
- [go-sqlite3](github.com/mattn/go-sqlite3) \- A small, fast SQL database implementation