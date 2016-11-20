package main

import (
	"flag"
	"github.com/dfree1645/nntp_linebot"
)

func main() {
	var (
		addr   = flag.String("addr", ":8080", "addr to bind")
		dbconf = flag.String("dbconf", "dbconfig.yml", "database configuration file.")
		env    = flag.String("env", "development", "application envirionment (production, development etc.)")
		conf   = flag.String("conf", "config.yml", "configuration file without db. (SSH user, nntp, line)")
		path   = flag.String("path", "./", "relative path to absolute path")
	)
	flag.Parse()
	b := base.New(*path)
	b.Init(*conf, *dbconf, *env, *path)
	b.Run(*addr)
}
