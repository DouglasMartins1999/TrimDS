package lib

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"path"
	"strings"

	"github.com/samber/lo"
)

var BUFFERSIZE = 1000000
var ROMTYPE = map[int8]string{0: "NDS", 1: "DSi", 2: "Enhanced"}

type ROM struct {
	Path     string
	FileSize int64
	ROMSize  int64
	Unit     int8
	WiFi     bool
}

func (r *ROM) CalcFileSize() *ROM {
	var file = lo.Must(os.Open(r.Path))
	var stat = lo.Must(file.Stat())

	file.Close()
	r.FileSize = stat.Size()
	return r
}

func (r *ROM) CalcROMSize() *ROM {
	var file = lo.Must(os.Open(r.Path))
	var unit int8
	var size uint32

	file.Seek(0x12, io.SeekStart)
	binary.Read(file, binary.LittleEndian, &unit)

	file.Seek(lo.Ternary[int64](unit == 0, 0x80, 0x210), io.SeekStart)
	binary.Read(file, binary.LittleEndian, &size)

	file.Close()

	r.ROMSize = int64(size)
	r.Unit = unit

	if r.ROMSize == 0 {
		r.ROMSize = r.FileSize
	}

	return r
}

func (r *ROM) CheckWiFi() *ROM {
	var input = lo.Must(os.Open(r.Path))
	var wifiLen int64 = 0x88

	if r.FileSize >= r.ROMSize+wifiLen {
		var wifi_data = make([]byte, wifiLen)
		var wifi_comp_00 = make([]byte, wifiLen)
		var wifi_comp_FF = bytes.Repeat([]byte{0xFF}, int(wifiLen))

		input.Seek(r.ROMSize, io.SeekStart)
		binary.Read(input, binary.LittleEndian, &wifi_data)

		if !bytes.Equal(wifi_data, wifi_comp_00) && !bytes.Equal(wifi_data, wifi_comp_FF) {
			r.WiFi = true
			r.ROMSize += wifiLen
		} else {
			r.WiFi = false
		}
	}

	input.Close()
	return r
}

func (r *ROM) Trim() bool {
	var input = lo.Must(os.Open(r.Path))
	var output = lo.Must(os.Create(r.trimPath()))

	defer input.Close()
	defer output.Close()

	if r.FileSize < 0x200 {
		return false
	}

	buf := make([]byte, BUFFERSIZE)
	input.Seek(0, io.SeekStart)
	copied := 0
	totalsize := int(r.ROMSize)

	for {
		n, err := input.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if copied+n > totalsize {
			n = totalsize - copied
		}

		if n == 0 {
			break
		}

		if _, err := output.Write(buf[:n]); err != nil {
			panic(err)
		}

		copied += n
	}

	return true
}

func readBytes(input *os.File, offset int64, data any) {
	input.Seek(0, io.SeekStart)
	input.Seek(offset, io.SeekStart)
	binary.Read(input, binary.LittleEndian, &data)
}

func (r *ROM) trimPath() string {
	oldBase := strings.TrimSuffix(path.Base(r.Path), ".nds")
	newPath := path.Join(path.Dir(r.Path), oldBase+".trim.nds")
	return newPath
}
