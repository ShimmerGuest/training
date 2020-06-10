package main

import (
	`encoding/binary`
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
	
	rs := RtmpSession{
		c:c,
	}
	err := receiveMsg(&rs)
	if err != nil {
	
	}
	
	//close 只作为警告处理
	err := c.Close()
	
	
}



func receiveMsg(s *RtmpSession) error {
	
	var err error
	var fmt byte
	bootstrap := make([]byte, 11)
	for {
		_, err := io.ReadAtLeast(s.c, bootstrap[:1], 1)
		if err != nil {
			return err
		}
	
		//forward 2 bit
		fmt = (bootstrap[0] >> 6)& 0x03
		//last 6 bit
		csid := int(bootstrap[0] & 0x3f)
		
		
		//ngx的read状态比较复杂，每个work只有一个事件循环，不能有任何阻塞操作
		//怎么做到chunk解析的，如果到对应的判断处发现buf不够。理没有读到的话，
		//下一次还要继续恢复现场后到这个，每次都是read完后都是需要保证下一次还能到这块。
		//ngx rtmp这块的代码可以作为一个通用的非阻塞c读取chunk的解决方案
		
		switch csid {
		case 0:
			//read 1 byte, csid
			if _, err = io.ReadAtLeast(s.c, bootstrap[:1], 1); err != nil {return err}
			csid = int(bootstrap[0]) + 64
			break
		case 1:
			if _, err = io.ReadAtLeast(s.c, bootstrap[:2], 2); err != nil {return err}
			csid = int(bootstrap[0]) + int(bootstrap[1]) * 256
			break
		//case 2:
		//	//common?
		//	break
		default:
			//noop
		}
		
		//将当前流的msg和head放到流中
		s := getStream(csid)
		s.head.csid = 1
		
		//开始解析Msg
		
		//switch fmt {
		//case 3:
		//case 2:
		//}
		if fmt <= 2 {
		
		
		}
	}
}

var csid2stream map[int]*Stream

func getStream(csid int) *Stream{
	
	s, ok := csid2stream[csid]
	if !ok {
		csid2stream[csid] = &Stream{
			head:   Header{
				csid:csid,
			},
			msgBuf: make([]byte, 100),
		}
	}
}

type Header struct {
	csid int
	timestamp int64
	msgType int
	msgId int
}

type Stream struct {
	head    Header
	msgBuf []byte
}

type RtmpSession struct {
	c net.Conn
}



