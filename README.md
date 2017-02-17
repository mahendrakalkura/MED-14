How to install?
===============

```
$ psql -c 'CREATE DATABASE "MED-14"' -d postgres
$ mkdir MED-14
$ cd MED-14
$ git clone --recursive git@github.com:mahendrakalkura/MED-14.git .
$ cp settings.toml.sample settings.toml
$ go get
```

How to run?
===========

```
$ cd MED-14
$ go build
$ ./MED-14 --action=bootstrap
$ ./MED-14 --action=insert
$ ./MED-14 --action=addresses_one
$ ./MED-14 --action=addresses_all
$ ./MED-14 --action=progress
$ ./MED-14 --action=report
```
