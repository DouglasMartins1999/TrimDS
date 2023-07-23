package lib

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"math"
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

func AddROM(path string) *ROM {
	if _, exists := FindROM(path); exists {
		return nil
	}

	rom := ROM{Path: path}
	rom.CalcFileSize().CalcROMSize().CheckWiFi()
	queue = append(queue, rom)
	return &rom
}

func AddROMFolder(folder string) []ROM {
	dirFiles := lo.Must(ioutil.ReadDir(folder))
	ndsRoms := lo.Filter(dirFiles, func(f fs.FileInfo, i int) bool {
		ext := path.Ext(f.Name())
		_, exists := FindROM(path.Join(folder, f.Name()))
		return !f.IsDir() && ext == ".nds" && !exists
	})

	newFiles := lo.Map(ndsRoms, func(f fs.FileInfo, i int) ROM {
		rom := ROM{Path: path.Join(folder, f.Name())}
		rom.CalcFileSize().CalcROMSize().CheckWiFi()
		return rom
	})

	queue = append(queue, newFiles...)
	return newFiles
}

func DeleteROM(index int32) {
	queue = append(queue[:index], queue[index+1:]...)
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
