package main

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/qx66/basicDiag/internal/biz"
	"net/url"
	"time"
	
	lTheme "github.com/qx66/basicDiag/pkg/theme"
)

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(&lTheme.MyTheme{})
	win := myApp.NewWindow("StartOps诊断系统")
	win.SetContent(makeUI(win))
	win.Resize(fyne.NewSize(win.Canvas().Size().Width, 650))
	win.ShowAndRun()
}

func makeUI(w fyne.Window) fyne.CanvasObject {
	ctx := context.Background()
	// header
	header := canvas.NewText("StartOps网络诊断", theme.PrimaryColor())
	header.TextSize = 42
	header.Alignment = fyne.TextAlignCenter
	
	// foot
	u, _ := url.Parse("https://startops.com.cn")
	footer := widget.NewHyperlinkWithStyle("startops.com.cn", u, fyne.TextAlignCenter, fyne.TextStyle{})
	
	//
	input := widget.NewEntry()
	input.MultiLine = true
	input.Wrapping = fyne.TextWrapBreak
	input.SetPlaceHolder("请输入需要诊断的链接")
	
	output := widget.NewEntry()
	output.MultiLine = true
	output.Wrapping = fyne.TextWrapBreak
	output.SetPlaceHolder("Output Result")
	
	diag := widget.NewButtonWithIcon(
		"诊断",
		theme.MediaSkipNextIcon(),
		func() {
			if input.Text == "" {
				//input.Text = w.Clipboard().Content()
				input.Refresh()
			}
			
			r, err := biz.BasicDiag(ctx, input.Text, output)
			if err != nil {
				output.Text = fmt.Sprintf("Time: %s, %s", time.Now().String(), err.Error())
			} else {
				output.Text = r
			}
			
			/*
				id, err := biz.BasicDiag(ctx, input.Text)
				if err != nil {
					output.Text = fmt.Sprintf("Time: %s, %s", time.Now().String(), err.Error())
				} else {
					output.Text = fmt.Sprintf("Time: %s, 诊断上报成功, id: %s", time.Now().String(), id)
				}
			*/
			
			output.Refresh()
		})
	
	diag.Importance = widget.HighImportance
	
	clear := widget.NewButtonWithIcon(
		"clear",
		theme.ContentClearIcon(),
		func() {
			output.Text = ""
			output.Refresh()
			input.Text = ""
			input.Refresh()
		},
	)
	clear.Importance = widget.MediumImportance
	
	/*
		decode := widget.NewButtonWithIcon("Decode", theme.MediaSkipPreviousIcon(), func() {
			if input.Text == "" {
				input.Text = w.Clipboard().Content()
				input.Refresh()
			}
			out, err := base64.StdEncoding.DecodeString(input.Text)
			if err == nil {
				output.Text = string(out)
			} else {
				output.Text = err.Error()
			}
			output.Text = string(out)
			output.Refresh()
		})
		decode.Importance = widget.HighImportance
	*/
	
	copy := widget.NewButtonWithIcon(
		"拷贝结果",
		theme.ContentCutIcon(),
		func() {
			clipboard := w.Clipboard()
			clipboard.SetContent(output.Text)
			//output.Text = ""
			//output.Refresh()
			//input.Text = ""
			//input.Refresh()
		},
	)
	copy.Importance = widget.WarningImportance
	
	return container.NewBorder(
		header,
		footer,
		nil,
		nil,
		container.NewGridWithRows(
			2,
			container.NewBorder(
				nil,
				//container.NewVBox(container.NewGridWithColumns(3, diag, clear, decode), copy),
				container.NewVBox(container.NewGridWithColumns(3, diag, clear, copy)),
				nil,
				nil,
				input),
			output,
		),
	)
}
