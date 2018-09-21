# Feed Me!

Feed Me! is a cross platform, open source RSS / Atom feed reader! Right now it
should work on Windows, Linux, and OSX.

## Building

### SQLite

Requires [sqlite3](https://www.sqlite.org/index.html) to run.

To get it on OSX via brew run

```
brew install sqlite3
```

On Ubuntu run

```
sudo apt-get install build-essential
```

and on Windows make sure the precompiled sqlite binaries are somewhere your
gcc toolchain can see them. Not sure if this works with MSVC but if you can
get sqlite to run on MSVC this should work as well.

### Go

You'll have to install [go](https://golang.org) then run

```
go get github.com/Noofbiz/FeedMe
```

Then to run you can just do

```
go run main.go
```

or if you want a compiled binary

```
go build main.go -o=FeedMe
./FeedMe
```
