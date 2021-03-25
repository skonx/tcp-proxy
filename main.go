package main

import (
	"io"
	"log"
	"net"
	"os"
)

const remote = "netflix.com:443"

func handle(s net.Conn) {
	d, err := net.Dial("tcp", remote)
	if err != nil {
		log.Fatalf("cannot connect to %s : %v\n", remote, err)
	}
	defer d.Close()

	dmw := io.MultiWriter(d, os.Stdout)
	smw := io.MultiWriter(s, os.Stdout)

	go func() {
		if _, err := io.Copy(dmw, s); err != nil {
			log.Fatalf("cannot copy data to %s : %v\n", remote, err)
		}
	}()

	if _, err := io.Copy(smw, d); err != nil {
		log.Fatalf("cannot copy data from %s : %v\n", remote, err)
	}

}

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("cannot start proxy server on port 80")
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatalln("cannot accept TCP connection")
		}
		go handle(c)
	}
}
