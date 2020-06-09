package main

import (
	`fmt`
	`io`
	`net`
)

func Run(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept connect fail")
			continue
		}
		
		go handlerConnect(c)
	}
}


func handlerConnect(c net.Conn) {
	//开始握手
	err := handshake(c)
	if err != nil {
		fmt.Println("handshake fail err = ", err.Error())
	}
	
	RecMsg()
}



func receiveMsg(c net.Conn) error {
	
	var err error
	var fmt byte
	
	bootstrap := make([]byte, 11)
	
	for {
		_, err := io.ReadAtLeast(c, bootstrap, 1)
		if err != nil {
			return err
		}
	
		//forward 2 bit
		fmt = bootstrap[0] & 0xc0
		
		//last 6 bit
		switch bootstrap[0] & 0x3f {
		//type 0-> next 1 bit
		case 0:
			//read 1 byte, csid
			_, err = io.ReadAtLeast(c, bootstrap, 1)
			break
		case 1:
			break
		case 2:
			break
		default:
		
		}
		
		
	}
}

type RtmpSession struct {
	CSID int
}



