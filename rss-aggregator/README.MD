### RSS Aggregator

A simple golang API server to parse RSS feeds, save post to database and server latest posts to user

### Endpoint

* /user                         POST    create user
* /user                         GET     get user (authentication required)

* /feeds                        POST   create a feed (authentication required)
* /feeds                        GET     get all feeds

* /feed_follows                 GET     get feed follows
* /feed_follows                 POST    follow a feed
* /feed_follows/{feedFollowID}  DELETE  unfollow a feed

* /posts                        GET     get posts for user

The server is configured to read the RSS feed urls from database, then scheduled to fech posts every minutes in a number of goroutines

### Techs

1. Golang
2. Postgres
3. [goose](https://github.com/pressly/goose) & [sqlc](https://github.com/sqlc-dev/sqlc)
4. Other 3rd party libs

### How to run

1. Clone the repo
2. Use terminal, go to the cloned folder
3. type: go mod vendor
4. type: go mod run .
5. Use postman or similar apps to explore the api

Note: you should have golang and postgres set up on your local machine, change the environment variables in .env if needed.