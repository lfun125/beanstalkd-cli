package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kr/beanstalk"
)

var (
	server string
	port   int
)

func main() {
	flag.StringVar(&server, "s", "127.0.0.1", "The server name where beanstalkd is running")
	flag.IntVar(&port, "p", 11300, "The port on which beanstalkd is listening")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		return
	}
	switch strings.ToUpper(args[0]) {
	case "PUT":
		put(args[1:]...)
	}
}

func put(args ...string) {
	if len(args) < 2 {
		fmt.Println("beanstalkd-cli put <tube> <data>")
		return
	}
	tubeName := args[0]
	data := args[1]
	var conn *beanstalk.Conn
	conn, err := beanstalk.Dial("tcp", fmt.Sprintf("%s:%d", server, port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	tube := beanstalk.Tube{Conn: conn, Name: tubeName}
	id, err := tube.Put([]byte(data), 1024, 0, time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("put success tube: %s, data: %s, id: %d \n", tubeName, data, id)
}
