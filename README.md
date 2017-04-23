# flasher
A golang implementation of the flash cards site [jwasham/computer-science-flash-cards](https://github.com/jwasham/computer-science-flash-cards) for practicing computer science concepts.

## Demo

Check the running application [here](https://flashergo.herokuapp.com).

```
Username: admin
Password: admin
```

## Built with
- Golang
- Gorilla Toolkit: gorilla/mux, gorilla/session, gorilla/context
- flosch/pongo2: a Django-syntax like templating-language
- SQLite

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) installed.

```sh
$ go get github.com/ympons/flasher
$ PORT=5800 $GOPATH/bin/flasher
```

Your local copy should now be running on [localhost:5800](http://localhost:5800/). Use `admin` as username and password to login to the website.
