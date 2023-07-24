package gui

import (
	"fmt"

	"dotins.eu.org/trimds/lib"
	"github.com/samber/lo"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

func RenderItem(index int, rom lib.ROM) {
	item := MainForm.FileList.Items().Add()
	item.SetCaption(fmt.Sprintf("%d", index))
	item.SubItems().Add(fmt.Sprintf(rom.Path))
	item.SubItems().Add(lib.FormatBytes(rom.FileSize, false))
	item.SubItems().Add(lib.FormatBytes(rom.ROMSize, false))
	item.SubItems().Add(lib.ROMTYPE[rom.Unit])
	item.SubItems().Add(lo.Ternary(rom.WiFi, "YES", "NO"))
	item.SetChecked(rom.FileSize < rom.ROMSize)
}

func RenderList() {
	files := lib.Values()

	MainForm.FileList.Clear()
	MainForm.FileList.Items().BeginUpdate()
	for i := 0; i < len(files); i++ {
		RenderItem(i+1, files[i])
	}
	MainForm.FileList.Items().EndUpdate()
}

func RenderStatus() {
	statusBar := MainForm.StatusBar.TStatusBar
	statusBar.Panels().Clear()

	if lib.Size() > 0 {
		totalFileSize := lib.TotalFileSize()
		totalRomSize := lib.TotalROMSize()
		totalSaved := 100 - (float64(totalRomSize) / float64(totalFileSize) * 100)

		filesize := lib.FormatBytes(totalFileSize, true)
		romsize := lib.FormatBytes(totalRomSize, true)
		saved := lib.FormatBytes(totalFileSize-totalRomSize, false)

		totalPanel := statusBar.Panels().Add()
		totalPanel.SetText(fmt.Sprintf("Total: %d", lib.Size()))
		totalPanel.SetWidth(120)

		sizePanel := statusBar.Panels().Add()
		sizePanel.SetText(fmt.Sprintf("Size: %s / Trimmed: %s - %s (%.1f%%)", filesize, romsize, saved, totalSaved))
		sizePanel.SetAlignment(types.TaLeftJustify)
	}
}

func AddROM(sender vcl.IObject) {
	if MainForm.DialogFile.Execute() {
		files := MainForm.DialogFile.Files()
		paths := []string{}

		for i := int32(0); i < files.Count(); i++ {
			paths = append(paths, files.S(i))
		}

		lib.AddROMs(paths)

		RenderList()
		RenderStatus()
	}
}

func AddROMFolder(sender vcl.IObject) {
	if MainForm.DialogFolder.Execute() {
		lib.AddROMFolder(MainForm.DialogFolder.FileName())
		RenderList()
		RenderStatus()
	}
}

func DeleteROM(sender vcl.IObject) {
	index := MainForm.FileList.ItemIndex()

	if index != -1 {
		lib.DeleteROM(int(index))
		RenderList()
		RenderStatus()
	}
}

func TrimROMQueue(sender vcl.IObject) {
	lib.Trim()
}
