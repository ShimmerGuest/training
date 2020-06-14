package main

import (
	`encoding/binary`
	`fmt`
	`io`
	`net`
)

const default_chunk_szie = 128

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
		chunkSize:default_chunk_szie,
	}
	
	err := receiveMsg(&rs)
	if err != nil {
	}
	
	//close 只作为警告处理
	err := c.Close()
}



func receiveMsg(rs *RtmpSession) error {
	
	var err error
	var fmt byte
	bootstrap := make([]byte, 11)
	for {
		_, err = io.ReadAtLeast(rs.c, bootstrap[:1], 1)
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
			if _, err = io.ReadAtLeast(rs.c, bootstrap[:1], 1); err != nil {return err}
			csid = int(bootstrap[0]) + 64
			break
		case 1:
			if _, err = io.ReadAtLeast(rs.c, bootstrap[:2], 2); err != nil {return err}
			csid = int(bootstrap[0]) + int(bootstrap[1]) * 256
			break
		//case 2:
		//	//common?
		//	break
		default:
			//noop
		}
		
		//将当前流的msg和head放到流中
		st := getStream(csid)
		st.m.head.csid = 1
		
		//开始解析Msg
		switch fmt {
		case 0:
			if _, err = io.ReadAtLeast(rs.c, bootstrap, 11); err != nil {return err}
			st.m.clear()
			st.m.head.timestamp =int64(BigEnd24(bootstrap))
			st.m.head.msgLen = int(BigEnd24(bootstrap[3:]))
			st.m.head.msgType = int(bootstrap[6])
			st.m.head.msId = binary.BigEndian.Uint32(bootstrap[7:])
		case 1:
			if _, err = io.ReadAtLeast(rs.c, bootstrap, 11); err != nil {return err}
			st.m.head.timestampDelta = int64(BigEnd24(bootstrap))
			st.m.head.timestamp += st.m.head.timestampDelta
			st.m.head.msgLen = int(BigEnd24(bootstrap[3:]))
			st.m.head.msgType = int(bootstrap[6])
		case 2:
			st.m.head.timestampDelta = int64(BigEnd24(bootstrap))
			st.m.head.timestamp += st.m.head.timestampDelta
		case 3:
		}
		
		//计算还需要多少byte拼接完这个message
		//chunk的大小都是固定的
		//最后一个chunk可能会小于chunk大小。
		
		//计算需要读取数据的大小
		needSize := st.m.head.msgLen - len(st.m.msgBuf)
		if needSize > rs.chunkSize {
			//如果大于一个chunk那么就需要读取一个chunk数据
			needSize = rs.chunkSize
		}
		
		chunkData := make([]byte, needSize)
		//
		if _, err = io.ReadAtLeast(rs.c, chunkData, needSize); err != nil {return err}
		
		//
		st.m.feed(chunkData)
		
		//当前已经满了
		if needSize > rs.chunkSize {
			 //handler msg
		}
	}
}

var csid2stream map[int]*Stream

func getStream(csid int) *Stream{
	
	s, ok := csid2stream[csid]
	if !ok {
		s = &Stream{
			m:Message{
				head:   Header{
					csid:           csid,
					msId:           0,
					timestamp:      0,
					timestampDelta: 0,
					msgType:        0,
					msgId:          0,
					msgLen:         0,
				},
				msgBuf: make([]byte, 100),
			},
		}
		csid2stream[csid] = s
	}
	
	return s
}

//头
type Header struct {
	csid int
	msId uint32
	timestamp int64
	timestampDelta int64
	msgType int
	msgId int
	msgLen int
}

//
type Stream struct {
	m Message
	
	//推流session
	pub RtmpSession
	
	//拉流session
	sub []RtmpSession
}


func (s *Stream) Msg()  *Message{
	return &s.m
}

type RtmpSession struct {
	c net.Conn
	chunkSize int
}

type Message struct {
	head Header
	msgBuf []byte
}

func (m* Message) feed(data []byte) {
	m.msgBuf = append(m.msgBuf, data...)
}

func (m* Message)clear() {
	m.msgBuf = m.msgBuf[:0]
}