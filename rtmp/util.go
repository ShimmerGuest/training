package main


func BigEnd24(bs []byte) int32 {
	return int32(bs[0]) << 16 | int32(bs[1]) << 8 | int32(bs[2])
}
