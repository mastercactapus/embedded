package main

import (
	"context"
	"log"
	"net"

	"github.com/mastercactapus/embedded/bustool"
)

func main() {
	n, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := n.Accept()
		if err != nil {
			panic(err)
		}

		sh := bustool.NewShell(conn, conn)
		bustool.AddI2C(sh, nil, nil)

		err = sh.Exec(context.Background())
		if err != nil {
			log.Println("ERROR:", err)
		}
		conn.Close()
	}
}
