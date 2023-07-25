package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"hash/crc32"
	"os"
)

func main() {
	genresByte(os.Args[1], os.Args[2])
}

func genresByte(source string, dest string) {
	input, _ := os.ReadFile(source)
	crc32Val := crc32.ChecksumIEEE(input)
	bs := compress(input)
	code := bytes.NewBuffer(nil)

	code.WriteString("package liblclbinres")
	code.WriteString("\r\n\r\n")
	code.WriteString(fmt.Sprintf("const CRC32Value uint32 = 0x%x\r\n\r\n", crc32Val))

	code.WriteString("var LCLBinRes = []byte(\"")
	for _, b := range bs {
		code.WriteString("\\x" + fmt.Sprintf("%.2x", b))
	}
	code.WriteString("\")\r\n")

	os.WriteFile(dest, code.Bytes(), 0666)
}

func compress(input []byte) []byte {
	var in bytes.Buffer

	w, _ := zlib.NewWriterLevel(&in, zlib.BestCompression)
	w.Write(input)
	w.Close()

	return in.Bytes()
}
