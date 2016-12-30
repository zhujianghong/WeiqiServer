// Tea.go
package Tea

import (
	"bytes"
	"encoding/binary"
)

var (
	teaKey uint32 = 0x26EDC334
	a      uint32 = teaKey
	b      uint32 = teaKey + 0x050E7F8D
	c      uint32 = teaKey + 0x10984F7E
	d      uint32 = teaKey + 0x76EF3720
	delta  uint32 = 0x9E3779B9
)

func Encrypt(data []byte, offset int, length int, times int) []byte {
	dataLength := length - offset
	in := bytes.NewBuffer(data[offset:length])
	out := make([]byte, dataLength)
	flag := 0
	cnt := dataLength / 8
	for i := 0; i < cnt; i++ {
		var (
			y   uint32 = 0
			z   uint32 = 0
			sum uint32 = 0
		)
		binary.Read(in, binary.LittleEndian, &y)
		binary.Read(in, binary.LittleEndian, &z)
		for j := 0; j < times; j++ {
			sum += delta
			y += ((z << 4) + a) ^ (z + sum) ^ ((z >> 5) + b)
			z += ((y << 4) + c) ^ (y + sum) ^ ((y >> 5) + d)
		}
		binary.LittleEndian.PutUint32(out[0+8*i:4+8*i], y)
		binary.LittleEndian.PutUint32(out[4+8*i:8+8*i], z)
		flag += 8
	}
	for i := 0; i < (dataLength - cnt*8); i++ {
		temp, _ := in.ReadByte()
		out[flag] = (temp ^ 0xff)
		flag++
	}
	return out
}

func Decrypt(data []byte, offset int, length int, times int) []byte {
	dataLength := length - offset
	in := bytes.NewBuffer(data[offset:length])
	out := make([]byte, dataLength)
	flag := 0
	cnt := dataLength / 8
	for i := 0; i < cnt; i++ {
		var (
			y   uint32 = 0
			z   uint32 = 0
			sum uint32 = 0xC6EF3720 //32
		)
		if times == 16 {
			sum = 0xE3779B90
		}
		binary.Read(in, binary.LittleEndian, &y)
		binary.Read(in, binary.LittleEndian, &z)
		for j := 0; j < times; j++ {
			z -= ((y << 4) + c) ^ (y + sum) ^ ((y >> 5) + d)
			y -= ((z << 4) + a) ^ (z + sum) ^ ((z >> 5) + b)
			sum -= delta
		}
		binary.LittleEndian.PutUint32(out[0+8*i:4+8*i], y)
		binary.LittleEndian.PutUint32(out[4+8*i:8+8*i], z)
		flag += 8
	}
	for i := 0; i < (dataLength - cnt*8); i++ {
		temp, _ := in.ReadByte()
		out[flag] = (temp ^ 0xff)
		flag++
	}
	return out
}
