package main

import (
	`bytes`
	`crypto/hmac`
	`crypto/sha256`
	`encoding/binary`
	`io`
	`media/src/library/errors`
	`net`
)

//30 + 32
var clientKey = []byte{
	'G', 'e', 'n', 'u', 'i', 'n', 'e', ' ', 'A', 'd', 'o', 'b', 'e', ' ',
	'F', 'l', 'a', 's', 'h', ' ', 'P', 'l', 'a', 'y', 'e', 'r', ' ',
	'0', '0', '1',
	
	0xF0, 0xEE, 0xC2, 0x4A, 0x80, 0x68, 0xBE, 0xE8, 0x2E, 0x00, 0xD0, 0xD1,
	0x02, 0x9E, 0x7E, 0x57, 0x6E, 0xEC, 0x5D, 0x2D, 0x29, 0x80, 0x6F, 0xAB,
	0x93, 0xB8, 0xE6, 0x36, 0xCF, 0xEB, 0x31, 0xAE,
}

// 36+32
var serverKey = []byte{
	'G', 'e', 'n', 'u', 'i', 'n', 'e', ' ', 'A', 'd', 'o', 'b', 'e', ' ',
	'F', 'l', 'a', 's', 'h', ' ', 'M', 'e', 'd', 'i', 'a', ' ',
	'S', 'e', 'r', 'v', 'e', 'r', ' ',
	'0', '0', '1',
	
	0xF0, 0xEE, 0xC2, 0x4A, 0x80, 0x68, 0xBE, 0xE8, 0x2E, 0x00, 0xD0, 0xD1,
	0x02, 0x9E, 0x7E, 0x57, 0x6E, 0xEC, 0x5D, 0x2D, 0x29, 0x80, 0x6F, 0xAB,
	0x93, 0xB8, 0xE6, 0x36, 0xCF, 0xEB, 0x31, 0xAE,
}

func handshake(c net.Conn) error {
	c0c1 := make([]byte, 1537)
	n, err := io.ReadFull(c, c0c1)
	if err != nil {
		return err
	}
	
	if n != 1537 {
		return errors.New("read message less than 1537")
	}
	
	if binary.BigEndian.Uint32(c0c1[5:]) == 0 {
		err = simpleHandshake(c, c0c1)
	} else {
		err = complexMode(c, c0c1)
	}
	
	return nil
}

func simpleHandshake(c net.Conn, c0c1 []byte) error {
	//
	s0s1s2 := make([]byte, 1537 + 1536)
	//s0s1s2[0] = c0c1[0]
	//
	//s1 := s0s1s2[1:]
	//binary.LittleEndian.PutUint32(s1, uint32(time.Now().Unix()))
	////写入0
	//binary.LittleEndian.PutUint32(s1, 0)
	////写入s1
	//c1 := c0c1[1:]
	//抄袭ngx把原样送回
	copy(s0s1s2, c0c1)
	s2 := s0s1s2[1537:]
	copy(s2, c0c1[1:])
	
	if n, err := c.Write(s2); err != nil{
		return err
	} else if n != len(s2) {
		return ErrWriteEnough
	}
	
	return nil
}


func complexMode(c net.Conn, c0c1 []byte) error {
	
	if findDigest(c0c1[1:], 8)  {
	
	}
	
	
	return nil
}


func findDigest(buf []byte, base int) bool {
	//也就是说这个offset是相对于自己的offset，并不是相对于c1的offset
	offset := (int(buf[base]) + int(buf[base+1]) + int(buf[base+2]) + int(buf[base+3])) % 728 + base + 4
	digest := make([]byte, 32)
	//对整个c1做digest
	makeDigest(buf, clientKey, digest, offset)
	
	if bytes.Compare(digest, buf[offset:offset + 32]) == 0 {
	
	}
	return
}

func makeDigest(b , key, out []byte, offset int) {
	mac := hmac.New(sha256.New, key)
	
	//left
	if offset != 0{
		mac.Write(b[:offset])
	}
	
	//right
	if offset + 32 < len(b) {
		mac.Write(b[offset+32:])
	}
	
	copy(out, mac.Sum(nil))
}

