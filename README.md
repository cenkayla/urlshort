# URL Shortener

- Written in Go
- Stores data in PostgreSQL with sqlx

# Quickstart
```shell
git clone https://github.com/cenkayla/shorturl.git
cd shorturl
go run main.go
```

# Usage
To create short url
```shell
$ curl -XPOST -d '{
"long_url":"example.com"
}' 'localhost:8080/create'
```
# Response
```shell
"U42wn"
```
