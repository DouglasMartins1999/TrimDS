package lib

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path"
	"runtime"

	"github.com/samber/lo"
)

var (
	queue    []ROM
	byteSize []string = []string{"bytes", "kB", "MB", "GB", "TB", "PB", "EB"}
)

func FindROM(path string) (ROM, bool) {
	return lo.Find(queue, func(r ROM) bool {
		return r.Path == path
	})
}

func AddROM(file string) *ROM {
	stat, _ := os.Stat(file)

	if stat.IsDir() {
		AddROMFolder(file)
	}

	if path.Ext(file) != ".nds" {
		return nil
	}

	if _, exists := FindROM(file); exists {
		return nil
	}

	rom := QueueROM(file)
	return &rom
}

func AddROMs(files []string) {
	for i := 0; i < len(files); i++ {
		AddROM(files[i])
	}
}

func AddROMFolder(folder string) []*ROM {
	dirFiles := lo.Must(os.ReadDir(folder))
	ndsRoms := lo.Filter(dirFiles, func(f fs.DirEntry, i int) bool {
		return !f.IsDir()
	})

	newFiles := lo.Map(ndsRoms, func(f fs.DirEntry, i int) *ROM {
		file := path.Join(folder, f.Name())
		return AddROM(file)
	})

	return newFiles
}

func QueueROM(file string) ROM {
	rom := ROM{Path: path.Clean(file)}
	rom.CalcFileSize().CalcROMSize().CheckWiFi()
	queue = append(queue, rom)
	return rom
}

func DeleteROM(index int) {
	queue = append(queue[:index], queue[index+1:]...)
}

func Clear() {
	queue = []ROM{}
}

func IndexOf(r ROM) int {
	return lo.IndexOf(queue, r)
}

func Size() int {
	return len(queue)
}

func Values() []ROM {
	return queue
}

func Trim() {
	for _, x := range Values() {
		if x.Trim() {
			DeleteROM(IndexOf(x))
		}
	}
}

func TotalFileSize() int64 {
	return lo.Reduce(queue, func(s int64, r ROM, i int) int64 {
		return s + r.FileSize
	}, 0)
}

func TotalROMSize() int64 {
	return lo.Reduce(queue, func(s int64, r ROM, i int) int64 {
		return s + r.ROMSize
	}, 0)
}

func FormatBytes(s int64, orig bool) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}

	base := lo.Ternary[float64](runtime.GOOS == "windows", 1024, 1000)
	log := math.Log(float64(s)) / math.Log(base)
	e := math.Floor(log)
	suffix := byteSize[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10

	if orig {
		return fmt.Sprintf("%.1f %s (%d bytes)", val, suffix, s)
	}

	return fmt.Sprintf("%.1f %s", val, suffix)
}
