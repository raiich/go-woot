package main

import (
	pb "../api"
	"../internal"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"strings"
	"time"
)

func parse() (port int, peers []string) {
	pointer := flag.Int("port", 8081, "port number")
	peer := flag.String("peer", "", "specify peer(s) by `host:port` with comma")
	flag.Parse()

	port = *pointer
	if *peer != "" {
		peers = strings.Split(*peer, ",")
	}
	return
}

func exec(peers []string) {
	var channels = make([]chan string, len(peers))

	println(len(peers))

	for i, peer := range peers {
		channels[i] = make(chan string, 1)
		go request(peer, channels[i])
	}
	for _, ch := range channels {
		go report(&ch)
	}
}

func request(peer string, ch chan string) {
	defer close(ch)
	conn, err := grpc.Dial(peer, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := pb.NewWootClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	body := fmt.Sprintf("%s, echo test!", peer)
	r, err := c.Echo(ctx, &pb.EchoRequest{Body: body})
	if err != nil {
		panic(err)
	}
	ch <- r.Body
}

func report(ch *chan string) {
	println("reporting")
	for s := range *ch {
		fmt.Println(s)
	}
}

func main() {
	port, peers := parse()
	fmt.Printf("starting server (port: %d)...\n", port)
	s := internal.StartServer(port)
	defer func() {
		time.Sleep(1 * time.Second)
		s.Stop()
		println("done...")
	}()

	exec(peers)
}
