# Potekhin Sergey & QBF

## Setup & run

```bash
$ sudo docker-compose up -d
$ sudo docker-compose logs --tail 100 app # Check that everything works fine
Attaching to qbf_test_app_1
app_1  | 10:39:31.871 DEBUG 001 Using database qbf.db
app_1  | 10:39:31.876 DEBUG 002 Using rate limiter with RPM 5
app_1  | 10:39:31.876 INFO 003 Starting server on 127.0.0.1:8080
$ sudo docker-compose exec app curl http://localhost:8080/ping # Test call
Pong
```

## TODO

- Integrate Sentry
- Move configuration to the environment
- GoDoc
- Docker fix
- Read at least 1 style guide for Go :/

## Endpoints

List of the endpoints, implemented in the API.

### Ping

Simple as always, just test the connection is fine.

```bash
curl -i http://localhost:8080/ping

HTTP/1.1 200 OK
Date: Tue, 02 Apr 2019 10:16:37 GMT
Content-Length: 4
Content-Type: text/plain; charset=utf-8

Pong
```

### Sync price

Get the price in a sync mode. Under the hood, the HTTP request will be performed to the Alpha Vantage API. Be carefull - in case you've reached the rate limit, you'll wait a lot to receive the response.

```bash
$ curl -i http://localhost:8080/price/sync?ticker=F

HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 02 Apr 2019 10:16:57 GMT
Content-Length: 38

{"price":"8.9800","ok":true,"msg":""}
```

### Async price

This endpoint allows to request the price without waiting for response. As soon as the price will be received, it will be saved into the database.

```bash
$ curl -i http://localhost:8080/price/async?ticker=F

HTTP/1.1 200 OK
Content-Type: application/json
Date: Tue, 02 Apr 2019 10:17:11 GMT
Content-Length: 54

{"price":"","ok":true,"msg":"Processing the request"}
```

### History

Returns the list of previous requests to the `price/sync` and `price/async`.

```bash
$ curl -i http://localhost:8080/history

HTTP/1.1 200 OK
Date: Tue, 02 Apr 2019 10:18:49 GMT
Content-Length: 1004
Content-Type: application/json

{"ok":true,"msg":"","history":[{"ID":1,"CreatedAt":"2019-04-02T11:26:23.403912155+03:00","UpdatedAt":"2019-04-02T11:26:23.403912155+03:00","DeletedAt":null,"Price":8.98}]}

```

### Health

Pretty simple endpoing - returns the list of records in the history requests list.

```bash
$ curl -i http://localhost:8080/health

HTTP/1.1 200 OK
Date: Tue, 02 Apr 2019 10:19:56 GMT
Content-Length: 33
Content-Type: text/plain; charset=utf-8

{"records":7,"ok":true,"msg":""}
```

## Walkthrough

```bash
$ tree .
.
├── docker-compose.yml
├── env
│   ├── dev.env
│   └── prod.env
├── main.go
├── qbf.db
├── README.md
└── vendor
    ├── app
    │   └── app.go
    ├── config
    │   └── config.go
    ├── connector
    │   └── connector.go
    ├── handlers
    │   ├── health.go
    │   ├── history.go
    │   └── price.go
    ├── models
    │   └── price.go
    └── utils
        └── utils.go
```

### Rate limiter

Really important part, implemented in a pure Go which is pretty cool. Allows you to be sure, that the rate limit of the data provider (in our case it's Alpha Vantage) will be considered. The main part of the rate limiter (so called `executor`) perform the following check each second:

1. Check that there're some pending requests
2. Check that the RPM still not reached (otherwise pass)
3. If it's possible - send the "allow" signal to the channel
4. Some service (e.g. `price/sync` handler) will be blocked by waiting for this signal and will continue the execution after receiving this signal.

## Used links

- [Task](https://www.evernote.com/shard/s495/client/snv?noteGuid=5ef58a18-b21f-4bf9-a3c0-cf9b3ba3527d&noteKey=fdc8144b577797a13e6c0906cc02b2d7&sn=https://www.evernote.com/shard/s495/sh/5ef58a18-b21f-4bf9-a3c0-cf9b3ba3527d/fdc8144b577797a13e6c0906cc02b2d7&title=%25D0%25A2%25D0%25B5%25D1%2581%25D1%2582%25D0%25BE%25D0%25B2%25D0%25BE%25D0%25B5%2B%25D0%25B7%25D0%25B0%25D0%25B4%25D0%25B0%25D0%25BD%25D0%25B8%25D0%25B5)
- [Structuring Applications in Go](https://medium.com/@benbjohnson/structuring-applications-in-go-3b04be4ff091)
- [Farewell Node.js](https://medium.com/@tjholowaychuk/farewell-node-js-4ba9e7f3e52b)
- [Alpha Vantage API Documentation](https://www.alphavantage.co/documentation/)
- [Введение в программирование на Go](http://golang-book.ru/)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
