package server

import (
	filedriver "github.com/goftp/file-driver"
	"github.com/goftp/server"
	"log"
)

func NewServer(conf *Config) *server.Server {
	factory := &filedriver.FileDriverFactory{
		RootPath: conf.Root,
		Perm:     server.NewSimplePerm("user", "group"),
	}

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     conf.Port,
		Hostname: conf.Host,
		Auth:     &server.SimpleAuth{Name: conf.User, Password: conf.Pass},
	}

	log.Printf("Starting ftp server on %v:%v", opts.Hostname, opts.Port)
	log.Printf("Username %v, Password %v", conf.User, conf.Pass)
	s := server.NewServer(opts)
	return s
}
