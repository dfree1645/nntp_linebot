DBNAME:=nntplinebot
ENV:=development

migrate/init:
	mysql -u root -h localhost --protocol tcp -e "create database `$(DBNAME)`" -p

migrate/up:
	sql-migrate up -env=$(ENV)

migrate/down:
	sql-migrate down -env=$(ENV)

migrate/status:
	sql-migrate status -env=$(ENV)

migrate/dry:
	sql-migrate up -dryrun -env=$(ENV)
