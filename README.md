JWT auth server
=============

Now supports only `mysql` driver

## Build

	make build

build `.deb` package

	make build-deb 		

## Deploy

Install `sql-migration` package

	go get -v github.com/rubenv/sql-migrate/...

Edit `dbconfig.yml` file 
	
Apply migrations from `./migrations`

	sql-migrate up
	
or for production env

	sql-migrate up -env=production
