package main

import (
	`net`
)

func Run(addr string) error {
	l, err := net.Listen("tcp", addr)
	
	for {
		c, err := l.Accept()
		go handlerConnect(c)
	}
}


func handlerConnect(c net.Conn) {
	//开始握手
	simpleHandshake(c)
}




