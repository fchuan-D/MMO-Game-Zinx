package main

import "lib/zinx/znet"

func main() {
	s := znet.NewServer()

	s.Serve()

}
