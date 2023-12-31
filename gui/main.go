package gui

import (
	"runtime"

	"dotins.eu.org/trimds/lib"
	"github.com/ying32/govcl/vcl"
	"github.com/ying32/govcl/vcl/types"
)

type TActionPanel struct {
	*vcl.TPanel

	AddFileBtn   *vcl.TButton
	AddFolderBtn *vcl.TButton
	TrimFilesBtn *vcl.TButton
}

type TStatusBar struct {
	*vcl.TStatusBar

	Total *vcl.TStatusPanel
	Size  *vcl.TStatusPanel
}

type TMainForm struct {
	*vcl.TForm

	MainMenu     *vcl.TMainMenu
	DialogFile   *vcl.TOpenDialog
	DialogFolder *vcl.TSelectDirectoryDialog
	FileList     *vcl.TListView
	StatusBar    TStatusBar
	ActionPanel  TActionPanel
}

var (
	MainForm *TMainForm
)

func Init(icon []byte) {
	if runtime.GOOS != "darwin" {
		appIcon := vcl.NewIcon()
		appIcon.LoadFromBytes(icon)
		vcl.Application.SetIcon(appIcon)
	}

	vcl.Application.Initialize()
	vcl.Application.SetMainFormOnTaskBar(true)
	vcl.Application.CreateForm(&MainForm)
	vcl.Application.Run()
}

func (f *TMainForm) OnFormCreate(sender vcl.IObject) {
	f.MainMenu = vcl.NewMainMenu(f)
	f.FileList = vcl.NewListView(f)
	f.DialogFile = vcl.NewOpenDialog(f)
	f.DialogFolder = vcl.NewSelectDirectoryDialog(f)
	f.StatusBar.TStatusBar = vcl.NewStatusBar(f)

	f.ActionPanel.TPanel = vcl.NewPanel(f)
	f.ActionPanel.AddFileBtn = vcl.NewButton(f)
	f.ActionPanel.AddFolderBtn = vcl.NewButton(f)
	f.ActionPanel.TrimFilesBtn = vcl.NewButton(f)

	f.FileList.SetParent(f)
	f.StatusBar.SetParent(f)
	f.ActionPanel.SetParent(f)

	configureMainForm(*MainForm.TForm)
	configureMainMenu(MainForm)
	configureDialogFile(*MainForm.DialogFile)
	configureDialogFolder(*MainForm.DialogFolder)
	configureListView(*MainForm.FileList)
	configureStatusBar(MainForm.StatusBar)
	configureActionPanel(MainForm.ActionPanel)
}

func (f *TMainForm) OnFormDropFiles(sender vcl.IObject, aFileNames []string) {
	lib.AddROMs(aFileNames)
	RenderList()
	RenderStatus()
}

func configureMainForm(f vcl.TForm) {
	f.SetCaption("TrimDS - Trim DS Roms")
	f.SetPosition(types.PoScreenCenter)

	f.SetWidth(700)
	f.SetHeight(400)
	f.Constraints().SetMaxHeight(600)
	f.Constraints().SetMaxWidth(1000)
	f.Constraints().SetMinHeight(300)
	f.Constraints().SetMinWidth(450)
	f.SetAllowDropFiles(true)
}
func configureMainMenu(f *TMainForm) {
	if runtime.GOOS == "darwin" {
		appMenu := vcl.NewMenuItem(f)

		appMenu.SetCaption(types.AppleLogoChar)
		subItem := vcl.NewMenuItem(f)

		subItem.SetCaption("About")
		subItem.SetOnClick(About)
		appMenu.Add(subItem)

		f.MainMenu.Items().Insert(0, appMenu)
	}

	item := vcl.NewMenuItem(f)
	item.SetCaption("File")

	subMenu := vcl.NewMenuItem(f)
	subMenu.SetCaption("Add File (&F)")
	subMenu.SetShortCutFromString("Ctrl+F")
	subMenu.SetOnClick(AddROM)
	item.Add(subMenu)

	subMenu = vcl.NewMenuItem(f)
	subMenu.SetCaption("Add Folder (&O)")
	subMenu.SetShortCutFromString("Ctrl+O")
	subMenu.SetOnClick(AddROMFolder)
	item.Add(subMenu)

	subMenu = vcl.NewMenuItem(f)
	subMenu.SetCaption("-")
	item.Add(subMenu)

	subMenu = vcl.NewMenuItem(f)
	subMenu.SetCaption("Trim Files (&T)")
	subMenu.SetShortCutFromString("Ctrl+T")
	subMenu.SetOnClick(TrimROMQueue)
	item.Add(subMenu)

	subMenu = vcl.NewMenuItem(f)
	subMenu.SetCaption("Clear List (&D)")
	subMenu.SetShortCutFromString("Ctrl+D")
	subMenu.SetOnClick(DeleteQueue)
	item.Add(subMenu)

	subMenu = vcl.NewMenuItem(f)
	subMenu.SetCaption("-")
	item.Add(subMenu)

	subMenu = vcl.NewMenuItem(f)
	subMenu.SetCaption("About (&A)")
	subMenu.SetShortCutFromString("Ctrl+A")
	subMenu.SetOnClick(About)
	item.Add(subMenu)

	subMenu = vcl.NewMenuItem(f)
	subMenu.SetCaption("Exit (&Q)")
	subMenu.SetShortCutFromString("Ctrl+Q")
	subMenu.SetOnClick(Exit)
	item.Add(subMenu)

	f.MainMenu.Items().Add(item)
}
func configureListView(lv vcl.TListView) {
	lv.SetAlign(types.AlClient)
	lv.SetRowSelect(false)
	lv.SetReadOnly(true)
	lv.SetViewStyle(types.VsReport)
	lv.SetGridLines(true)

	col := lv.Columns().Add()
	col.SetCaption("#")
	col.SetWidth(35)

	col = lv.Columns().Add()
	col.SetCaption("Name")
	col.SetAutoSize(true)
	col.SetMaxWidth(400)
	col.SetAlignment(types.MbRight)

	col = lv.Columns().Add()
	col.SetCaption("Size")
	col.SetWidth(100)
	col.SetAlignment(types.MbMiddle)

	col = lv.Columns().Add()
	col.SetCaption("Trimmed")
	col.SetWidth(100)
	col.SetAlignment(types.MbMiddle)

	col = lv.Columns().Add()
	col.SetCaption("ROM Type")
	col.SetWidth(100)
	col.SetAlignment(types.MbMiddle)

	col = lv.Columns().Add()
	col.SetCaption("Wireless")
	col.SetWidth(100)
	col.SetAlignment(types.MbMiddle)

	lv.SetOnDblClick(DeleteROM)
}
func configureDialogFile(d vcl.TOpenDialog) {
	d.SetFilter("NDS ROMS(*.nds)|*.nds")
	d.SetOptions(d.Options().Include(types.OfAllowMultiSelect))
	d.SetTitle("Choose NDS File")
}
func configureDialogFolder(d vcl.TSelectDirectoryDialog) {
	d.SetTitle("Choose Folder")
	d.SetOptions(d.Options().Include(types.OfViewDetail))
}
func configureStatusBar(s TStatusBar) {
	s.SetSimplePanel(false)
}
func configureActionPanel(pnl TActionPanel) {
	pnl.SetAlign(types.AlBottom)

	pnl.TrimFilesBtn.SetParent(pnl.TPanel)
	pnl.TrimFilesBtn.SetAnchors(types.NewSet(types.AkRight, types.AkBottom))
	pnl.TrimFilesBtn.SetTop(10)
	pnl.TrimFilesBtn.SetLeft(30)
	pnl.TrimFilesBtn.SetWidth(120)
	pnl.TrimFilesBtn.SetDefault(true)
	pnl.TrimFilesBtn.SetCaption("Trim Files")
	pnl.TrimFilesBtn.SetOnClick(TrimROMQueue)

	pnl.AddFileBtn.SetParent(pnl.TPanel)
	pnl.AddFileBtn.SetCaption("Add Files")
	pnl.AddFileBtn.SetWidth(120)
	pnl.AddFileBtn.SetTop(10)
	pnl.AddFileBtn.SetLeft(10)
	pnl.AddFileBtn.SetOnClick(AddROM)

	pnl.AddFolderBtn.SetParent(pnl.TPanel)
	pnl.AddFolderBtn.SetTop(10)
	pnl.AddFolderBtn.SetLeft(pnl.AddFileBtn.Left() + pnl.AddFileBtn.Width() + 10)
	pnl.AddFolderBtn.SetWidth(120)
	pnl.AddFolderBtn.SetCaption("Add Folder")
	pnl.AddFolderBtn.SetOnClick(AddROMFolder)
}
