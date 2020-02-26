package main

import (
	pb "../api"
	w "../internal"
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

func main0() {
	port, peers := parse()
	fmt.Printf("starting server (port: %d)...\n", port)
	s := w.StartServer(port)
	defer func() {
		time.Sleep(1 * time.Second)
		s.Stop()
		println("done...")
	}()
	exec(peers)
}

func eq(a *pb.Wid, b *pb.Wid) bool {
	return w.Equal(a, *b)
}

func clone(ss *w.SubSeq) []*pb.Wchar {
	ret := make([]*pb.Wchar, 0)
	for ; ss != nil; ss = ss.Next() {
		ret = append(ret, ss.Val())
	}
	return ret
}

func diff(former []*pb.Wchar, latter []*pb.Wchar) string {
	ret := make([]rune, 0)
	i, j := 0, 0
	for ; i < len(former) && j < len(latter); {
		o, n := former[i], latter[j]
		if eq(o.Id, n.Id) {
			if o.Visible {
				if n.Visible { // keep
					ret = append(ret, n.CodePoint)
				} else { // deleted
					ret = append(ret, ' ')
				}
			} else {
				if n.Visible {
					panic("invalid")
				}
			}
			i += 1
			j += 1
		} else {
			if !n.Visible {
				panic("invalid")
			} else { // inserted
				ret = append(ret, '?')
			}
		}
	}
	return string(ret)
}

func main() {
	site1 := w.NewSite("site1", "a")
	former := clone(site1.Raw())
	print("\033[H\033[2J")
	println(site1.Value())

	site1.GenerateIns(0, 'A')
	latter := clone(site1.Raw())
	time.Sleep(time.Second)
	print("\033[H\033[2J")
	println(diff(former, latter))
	println(site1.Value())

	site1.GenerateIns(2, 'B')
	former, latter = latter, clone(site1.Raw())
	time.Sleep(time.Second)
	print("\033[H\033[2J")
	println(diff(former, latter))
	println(site1.Value())

	site1.GenerateDel(1)
	former, latter = latter, clone(site1.Raw())
	time.Sleep(time.Second)
	print("\033[H\033[2J")
	println(diff(former, latter))
	println(site1.Value())
}
