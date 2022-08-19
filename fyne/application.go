package main

import (
	"fmt"
	"image/color"
	"log"
	"net/url"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

/*
Next:
Fix Scaling
Make bottomButtons work
*/

func main() { // rage in the darkness
	// we can add styles later

	var filePath string

	myApp := app.New()
	myWindow := myApp.NewWindow("ReTorrent")
	myWindow.Resize(fyne.NewSize(1000, 500))
	//	myWindow.SetFixedSize(true)

	appIcon, _ := fyne.LoadResourceFromURLString("https://static.wikia.nocookie.net/leagueoflegends/images/4/4e/Avarosan_profileicon.png")
	myWindow.SetIcon(appIcon)

	item1 := fyne.NewMenuItem("Open Torrent", func() {

		fileReader := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {

			if err != nil {
				dialog.ShowError(err, myWindow)
				return
			}
			if uc == nil {
				log.Println("Cancelled")
				return
			}

			filePath = uc.URI().String()
			fmt.Println(filePath) // gets the filePATH to the torrent

		}, myWindow)

		log.Println(myWindow.Canvas().Size().Width)

		fileReader.Resize(fyne.NewSize(500, 400))
		fileReader.Show()
	})

	item3 := fyne.NewMenuItem("Magent Link", func() {

	})

	item2 := fyne.NewMenuItem("About", func() {

		link, err := url.Parse("https://www.youtube.com/watch?v=2JpkMXinO1M&list=RDSZOr9K01Eno")

		if err == nil {
			myApp.OpenURL(link)

		}

	})
	header, item4 := bodyContainer(myWindow)
	torrentlist1, item5 := torrentList(myWindow)
	torrentlist2, item52 := torrentList2(myWindow)
	torrentlist3, item123 := torrentList3(myWindow)
	torrentlist4, item129 := torrentList4(myWindow)

	menu := fyne.NewMenu("File", item1, item3)

	menu2 := fyne.NewMenu("About", &item4, &item5, &item123, &item129, &item52, item2)

	myWindow.SetMainMenu(fyne.NewMainMenu(menu, menu2))

	// menu Done
	seperator1 := widget.NewSeparator()
	seperator1.Resize(fyne.NewSize(20, 0))

	split := container.NewVSplit(container.NewBorder(container.NewVBox(header, seperator1), nil, nil, nil, container.NewMax(canvas.NewRectangle(color.NRGBA{R: 19, G: 25, B: 41, A: 255}), container.NewVBox(torrentlist1, widget.NewSeparator(), torrentlist2, widget.NewSeparator(), torrentlist3, widget.NewSeparator(), torrentlist4))), bottomInfo(myWindow))

	split.Offset = 1.0
	myWindow.SetContent(split)

	log.Println("Hreres", myWindow.Canvas().Size())

	myWindow.ShowAndRun()

}

func torrentList2(win fyne.Window) (fyne.CanvasObject, fyne.MenuItem) {

	sty := &fyne.TextStyle{
		Bold: true,
	}

	item := widget.NewAccordionItem("Assignment Outer Space", widget.NewLabel("The Episode"))
	fileName := widget.NewAccordion(item)
	size := widget.NewLabelWithStyle("1.07GB", fyne.TextAlignCenter, *sty)
	proBar := widget.NewProgressBar()
	seeders := widget.NewLabelWithStyle("2", fyne.TextAlignCenter, *sty)
	leechers := widget.NewLabelWithStyle("9", fyne.TextAlignCenter, *sty)
	downloadSpeed := widget.NewLabelWithStyle("0.0 MB/s", fyne.TextAlignCenter, *sty)
	ETA := widget.NewLabelWithStyle("Completed ", fyne.TextAlignCenter, *sty)
	dateOfAdd := widget.NewLabelWithStyle("5/28/2022", fyne.TextAlignCenter, *sty)

	//ResizeAndMove(fileName, 0, 0, 298, 40)

	ResizeAndMove(size, 295, 0, 62, 40)

	ResizeAndMove(proBar, 380, 0, 120, 35)

	ResizeAndMove(seeders, 511, 0, 80, 40)

	ResizeAndMove(leechers, 592, 0, 80, 40)

	ResizeAndMove(downloadSpeed, 672, 0, 160, 40)

	ResizeAndMove(ETA, 832, 0, 80, 40)

	ResizeAndMove(dateOfAdd, 912, 0, 80, 40)
	proBar.SetValue(1)

	item3 := fyne.NewMenuItem("Adjust Window", func() {
		fmt.Println("Tap Size")
		lenghtDif := (win.Canvas().Size().Width - 1000) / 8
		fmt.Println("Tap CHange", lenghtDif)
		ResizeAndMove(fileName, 0, 0, 298+lenghtDif, 40)

		ResizeAndMove(size, 301+lenghtDif, 0, 62+lenghtDif, 40)

		ResizeAndMove(proBar, 365+(2*lenghtDif), 0, 145+(lenghtDif), 40)

		ResizeAndMove(seeders, 511+(3*lenghtDif), 0, 80+(lenghtDif), 40)

		ResizeAndMove(leechers, 592+(4*lenghtDif), 0, 80+(lenghtDif), 40)

		ResizeAndMove(downloadSpeed, 672+(5*lenghtDif), 0, 160+(lenghtDif), 40)

		ResizeAndMove(ETA, 832+(6*lenghtDif), 0, 80+(lenghtDif), 40)

		ResizeAndMove(dateOfAdd, 912+(7*lenghtDif), 0, 80+(lenghtDif), 40)

		win.Canvas().Content().Refresh()
		log.Println("Refresh:", win.Canvas().Size().Width)

	})

	list := widget.NewList(
		// lets change item count from 3 to 30
		func() int { return 1 }, // my list contain 3 items

		func() fyne.CanvasObject {

			//pro.Resize(fyne.NewSize(200, 50))

			//item1 := widget.NewAccordionItem("Attack On Titan S04Ep28", widget.NewLabel("The Episode"))

			//row := container.NewHBox(ac, canvas.NewRectangle(color.Black), size, pro, widget.NewLabel("11"), widget.NewLabel("18"), widget.NewLabel("5.4 MB/s"), widget.NewLabel("9 Min "), widget.NewLabel("6/2/2022"))
			//	container.NewAdaptiveGrid()

			row2 := container.NewWithoutLayout(fileName, size, proBar, seeders, leechers, downloadSpeed, ETA, dateOfAdd)

			return row2

		},
		// last one
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			// update data of widget
			//co.(*widget.Label).SetText("Here is my Newtext")

			//co.Resize(fyne.NewSize(150, 50))
		},
	)
	return list, *item3
}

func torrentList3(win fyne.Window) (fyne.CanvasObject, fyne.MenuItem) {

	sty := &fyne.TextStyle{
		Bold: true,
	}

	item := widget.NewAccordionItem("The Arts of flowers", widget.NewLabel("The Episode"))
	fileName := widget.NewAccordion(item)
	size := widget.NewLabelWithStyle("0.58GB", fyne.TextAlignCenter, *sty)
	proBar := widget.NewProgressBar()
	seeders := widget.NewLabelWithStyle("12", fyne.TextAlignCenter, *sty)
	leechers := widget.NewLabelWithStyle("25", fyne.TextAlignCenter, *sty)
	downloadSpeed := widget.NewLabelWithStyle("0.0 MB/s", fyne.TextAlignCenter, *sty)
	ETA := widget.NewLabelWithStyle("Completed ", fyne.TextAlignCenter, *sty)
	dateOfAdd := widget.NewLabelWithStyle("5/25/2022", fyne.TextAlignCenter, *sty)

	//ResizeAndMove(fileName, 0, 0, 298, 40)

	ResizeAndMove(size, 295, 0, 62, 40)

	ResizeAndMove(proBar, 380, 0, 120, 35)

	ResizeAndMove(seeders, 511, 0, 80, 40)

	ResizeAndMove(leechers, 592, 0, 80, 40)

	ResizeAndMove(downloadSpeed, 672, 0, 160, 40)

	ResizeAndMove(ETA, 832, 0, 80, 40)

	ResizeAndMove(dateOfAdd, 912, 0, 80, 40)
	proBar.SetValue(1)

	item3 := fyne.NewMenuItem("Adjust Window", func() {
		fmt.Println("Tap Size")
		lenghtDif := (win.Canvas().Size().Width - 1000) / 8
		fmt.Println("Tap CHange", lenghtDif)
		ResizeAndMove(fileName, 0, 0, 298+lenghtDif, 40)

		ResizeAndMove(size, 301+lenghtDif, 0, 62+lenghtDif, 40)

		ResizeAndMove(proBar, 365+(2*lenghtDif), 0, 145+(lenghtDif), 40)

		ResizeAndMove(seeders, 511+(3*lenghtDif), 0, 80+(lenghtDif), 40)

		ResizeAndMove(leechers, 592+(4*lenghtDif), 0, 80+(lenghtDif), 40)

		ResizeAndMove(downloadSpeed, 672+(5*lenghtDif), 0, 160+(lenghtDif), 40)

		ResizeAndMove(ETA, 832+(6*lenghtDif), 0, 80+(lenghtDif), 40)

		ResizeAndMove(dateOfAdd, 912+(7*lenghtDif), 0, 80+(lenghtDif), 40)

		win.Canvas().Content().Refresh()
		log.Println("Refresh:", win.Canvas().Size().Width)

	})

	list := widget.NewList(
		// lets change item count from 3 to 30
		func() int { return 1 }, // my list contain 3 items

		func() fyne.CanvasObject {

			//pro.Resize(fyne.NewSize(200, 50))

			//item1 := widget.NewAccordionItem("Attack On Titan S04Ep28", widget.NewLabel("The Episode"))

			//row := container.NewHBox(ac, canvas.NewRectangle(color.Black), size, pro, widget.NewLabel("11"), widget.NewLabel("18"), widget.NewLabel("5.4 MB/s"), widget.NewLabel("9 Min "), widget.NewLabel("6/2/2022"))
			//	container.NewAdaptiveGrid()

			row2 := container.NewWithoutLayout(fileName, size, proBar, seeders, leechers, downloadSpeed, ETA, dateOfAdd)

			return row2

		},
		// last one
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			// update data of widget
			//co.(*widget.Label).SetText("Here is my Newtext")

			//co.Resize(fyne.NewSize(150, 50))
		},
	)
	return list, *item3
}

func torrentList4(win fyne.Window) (fyne.CanvasObject, fyne.MenuItem) {

	sty := &fyne.TextStyle{
		Bold: true,
	}

	item := widget.NewAccordionItem("Wild Oceans ", widget.NewLabel("The Episode"))
	fileName := widget.NewAccordion(item)
	size := widget.NewLabelWithStyle("1.71 GB", fyne.TextAlignCenter, *sty)
	proBar := widget.NewProgressBar()
	seeders := widget.NewLabelWithStyle("50", fyne.TextAlignCenter, *sty)
	leechers := widget.NewLabelWithStyle("12", fyne.TextAlignCenter, *sty)
	downloadSpeed := widget.NewLabelWithStyle("0.0 MB/s", fyne.TextAlignCenter, *sty)
	ETA := widget.NewLabelWithStyle("Completed ", fyne.TextAlignCenter, *sty)
	dateOfAdd := widget.NewLabelWithStyle("6/1/2022", fyne.TextAlignCenter, *sty)

	//ResizeAndMove(fileName, 0, 0, 298, 40)

	ResizeAndMove(size, 295, 0, 62, 40)

	ResizeAndMove(proBar, 380, 0, 120, 35)

	ResizeAndMove(seeders, 511, 0, 80, 40)

	ResizeAndMove(leechers, 592, 0, 80, 40)

	ResizeAndMove(downloadSpeed, 672, 0, 160, 40)

	ResizeAndMove(ETA, 832, 0, 80, 40)

	ResizeAndMove(dateOfAdd, 912, 0, 80, 40)
	proBar.SetValue(1)

	item3 := fyne.NewMenuItem("Adjust Window", func() {
		fmt.Println("Tap Size")
		lenghtDif := (win.Canvas().Size().Width - 1000) / 8
		fmt.Println("Tap CHange", lenghtDif)
		ResizeAndMove(fileName, 0, 0, 298+lenghtDif, 40)

		ResizeAndMove(size, 301+lenghtDif, 0, 62+lenghtDif, 40)

		ResizeAndMove(proBar, 365+(2*lenghtDif), 0, 145+(lenghtDif), 40)

		ResizeAndMove(seeders, 511+(3*lenghtDif), 0, 80+(lenghtDif), 40)

		ResizeAndMove(leechers, 592+(4*lenghtDif), 0, 80+(lenghtDif), 40)

		ResizeAndMove(downloadSpeed, 672+(5*lenghtDif), 0, 160+(lenghtDif), 40)

		ResizeAndMove(ETA, 832+(6*lenghtDif), 0, 80+(lenghtDif), 40)

		ResizeAndMove(dateOfAdd, 912+(7*lenghtDif), 0, 80+(lenghtDif), 40)

		win.Canvas().Content().Refresh()
		log.Println("Refresh:", win.Canvas().Size().Width)

	})

	list := widget.NewList(
		// lets change item count from 3 to 30
		func() int { return 1 }, // my list contain 3 items

		func() fyne.CanvasObject {

			//pro.Resize(fyne.NewSize(200, 50))

			//item1 := widget.NewAccordionItem("Attack On Titan S04Ep28", widget.NewLabel("The Episode"))

			//row := container.NewHBox(ac, canvas.NewRectangle(color.Black), size, pro, widget.NewLabel("11"), widget.NewLabel("18"), widget.NewLabel("5.4 MB/s"), widget.NewLabel("9 Min "), widget.NewLabel("6/2/2022"))
			//	container.NewAdaptiveGrid()

			row2 := container.NewWithoutLayout(fileName, size, proBar, seeders, leechers, downloadSpeed, ETA, dateOfAdd)

			return row2

		},
		// last one
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			// update data of widget
			//co.(*widget.Label).SetText("Here is my Newtext")

			//co.Resize(fyne.NewSize(150, 50))
		},
	)
	return list, *item3
}

func bodyContainer(win fyne.Window) (fyne.CanvasObject, fyne.MenuItem) {

	log.Println("SIZE:", win.Canvas().Size().Width)

	content := container.NewMax()
	size := &widget.Button{
		Alignment: widget.ButtonAlignCenter,
		Text:      " Size ",
		OnTapped: func() {
			fmt.Println("Tap Size")

			log.Println("Refresh:", win.Canvas().Size().Width)

		},
		Importance: widget.HighImportance,
	}

	fileName := &widget.Button{
		Alignment:  widget.ButtonAlignCenter,
		Text:       "   File Name   ",
		OnTapped:   func() { fmt.Println("Tap") },
		Importance: widget.HighImportance,
	}

	progressBar := &widget.Button{
		Alignment:  widget.ButtonAlignCenter,
		Text:       " Progress ",
		OnTapped:   func() { fmt.Println("Tap Progress") },
		Importance: widget.HighImportance,
	}

	Seeders := &widget.Button{
		Alignment:  widget.ButtonAlignCenter,
		Text:       " Seeders ",
		OnTapped:   func() { fmt.Println("Tap Seeders") },
		Importance: widget.HighImportance,
	}

	Leechers := &widget.Button{
		Alignment:  widget.ButtonAlignCenter,
		Text:       " Leechers ",
		OnTapped:   func() { fmt.Println("Tap Seeders") },
		Importance: widget.HighImportance,
	}

	downloadSpeed := &widget.Button{
		Alignment:  widget.ButtonAlignCenter,
		Text:       " Download Speed ",
		OnTapped:   func() { fmt.Println("Tap Download") },
		Importance: widget.HighImportance,
	}

	ETA := &widget.Button{
		Alignment:  widget.ButtonAlignCenter,
		Text:       " ETA ",
		OnTapped:   func() { fmt.Println("Tap Download") },
		Importance: widget.HighImportance,
	}

	date := &widget.Button{
		Alignment:  widget.ButtonAlignCenter,
		Text:       " Date ",
		OnTapped:   func() { fmt.Println("Tap Download", ETA.Position()) },
		Importance: widget.HighImportance,
	}
	/*progressBar := &canvas.Text{
		//Alignment: fyne.TextAlignCenter,
		Text:     "Progress",
		TextSize: 17,
		Color:    color.White, //change color
		TextStyle: fyne.TextStyle{
			Bold: true,
		},
	}*/
	//progressBar := canvas.NewText("Progress", color.White)
	seperator := widget.NewSeparator()
	seperator2 := widget.NewSeparator()
	seperator3 := widget.NewSeparator()
	seperator4 := widget.NewSeparator()
	seperator5 := widget.NewSeparator()
	seperator6 := widget.NewSeparator()
	seperator7 := widget.NewSeparator()
	seperator8 := widget.NewSeparator()

	// 	treNodeIdd := []string{"widget.tre"}

	// //	treeTest := widget.NewTree(func(tni widget.TreeNodeID) []widget.TreeNodeID { return treNodeIdd }, func(tni widget.TreeNodeID) bool { return true }, func(b bool) fyne.CanvasObject {

	// 	return widget.NewButton("TreeButton", func() { fmt.Println("TreeButton") })
	// 	}, func(tni widget.TreeNodeID, b bool, co fyne.CanvasObject) {})
	//vertical := container.NewVBox(widget.NewLabel("test5"), widget.NewLabel("test3"), torrentList())
	item3 := fyne.NewMenuItem("Adjust Window", func() {
		fmt.Println("Tap Size")
		lenghtDif := (win.Canvas().Size().Width - 1000) / 8
		fmt.Println("Tap CHange", lenghtDif)
		ResizeAndMove(fileName, 0, 0, 298+lenghtDif, 40)
		ResizeAndMove(seperator, 298+lenghtDif, 0, 3, 40)

		ResizeAndMove(size, 301+lenghtDif, 0, 62+lenghtDif, 40)
		ResizeAndMove(seperator2, 362+(2*lenghtDif), 0, 3, 40)

		ResizeAndMove(progressBar, 365+(2*lenghtDif), 0, 145+(lenghtDif), 40)
		ResizeAndMove(seperator3, 508+(3*lenghtDif), 0, 3, 39)

		ResizeAndMove(Seeders, 511+(3*lenghtDif), 0, 80+(lenghtDif), 40)
		ResizeAndMove(seperator4, 590+(4*lenghtDif), 0, 3, 40)

		ResizeAndMove(Leechers, 592+(4*lenghtDif), 0, 80+(lenghtDif), 40)
		ResizeAndMove(seperator5, 671+(5*lenghtDif), 0, 3, 40)

		ResizeAndMove(downloadSpeed, 672+(5*lenghtDif), 0, 160+(lenghtDif), 40)
		ResizeAndMove(seperator6, 830+(6*lenghtDif), 0, 3, 40)

		ResizeAndMove(ETA, 832+(6*lenghtDif), 0, 80+(lenghtDif), 40)
		ResizeAndMove(seperator7, 910+(7*lenghtDif), 0, 3, 40)

		ResizeAndMove(date, 912+(7*lenghtDif), 0, 80+(lenghtDif), 40)
		ResizeAndMove(seperator8, 991+(8*lenghtDif), 0, 3, 40)

		win.Canvas().Content().Refresh()
		log.Println("Refresh:", win.Canvas().Size().Width)

	})

	log.Println("GONE HERE")

	ResizeAndMove(fileName, 0, 0, 298, 40)
	ResizeAndMove(seperator, 298, 0, 3, 40)

	ResizeAndMove(size, 301, 0, 62, 40)
	ResizeAndMove(seperator2, 362, 0, 3, 40)

	ResizeAndMove(progressBar, 365, 0, 145, 40)
	ResizeAndMove(seperator3, 508, 0, 3, 39)

	ResizeAndMove(Seeders, 511, 0, 80, 40)
	ResizeAndMove(seperator4, 590, 0, 3, 40)

	ResizeAndMove(Leechers, 592, 0, 80, 40)
	ResizeAndMove(seperator5, 671, 0, 3, 40)

	ResizeAndMove(downloadSpeed, 672, 0, 160, 40)
	ResizeAndMove(seperator6, 830, 0, 3, 40)

	ResizeAndMove(ETA, 832, 0, 80, 40)
	ResizeAndMove(seperator7, 910, 0, 3, 40)

	ResizeAndMove(date, 912, 0, 80, 40)
	ResizeAndMove(seperator8, 991, 0, 3, 40)

	header := container.NewWithoutLayout(
		fileName,
		seperator,
		size, seperator2, progressBar, seperator3, Seeders, seperator4, Leechers, seperator5, downloadSpeed, seperator6, ETA, seperator7, date, seperator8,
		content)

	return header, *item3

}
func ResizeAndMove(element fyne.CanvasObject, xPos, yPos, width, height float32) {
	element.Resize(fyne.NewSize(width, height))
	element.Move(fyne.NewPos(xPos, yPos))
}

func torrentList(win fyne.Window) (fyne.CanvasObject, fyne.MenuItem) {

	sty := &fyne.TextStyle{
		Bold: true,
	}

	item := widget.NewAccordionItem("Linux Mint 21 Cinnamon 64Bit ISO", widget.NewLabel("The Episode"))
	fileName := widget.NewAccordion(item)
	size := widget.NewLabelWithStyle("2.28GB", fyne.TextAlignCenter, *sty)
	proBar := widget.NewProgressBar()
	seeders := widget.NewLabelWithStyle("864", fyne.TextAlignCenter, *sty)
	leechers := widget.NewLabelWithStyle("24", fyne.TextAlignCenter, *sty)
	downloadSpeed := widget.NewLabelWithStyle("5.4 MB/s", fyne.TextAlignCenter, *sty)
	ETA := widget.NewLabelWithStyle("9 Min ", fyne.TextAlignCenter, *sty)
	dateOfAdd := widget.NewLabelWithStyle("6/2/2022", fyne.TextAlignCenter, *sty)

	//ResizeAndMove(fileName, 0, 0, 298, 40)

	ResizeAndMove(size, 295, 0, 62, 40)

	ResizeAndMove(proBar, 380, 0, 120, 35)

	ResizeAndMove(seeders, 511, 0, 80, 40)

	ResizeAndMove(leechers, 592, 0, 80, 40)

	ResizeAndMove(downloadSpeed, 672, 0, 160, 40)

	ResizeAndMove(ETA, 832, 0, 80, 40)

	ResizeAndMove(dateOfAdd, 912, 0, 80, 40)
	proBar.SetValue(.77)

	item3 := fyne.NewMenuItem("Adjust Window", func() {
		fmt.Println("Tap Size")
		lenghtDif := (win.Canvas().Size().Width - 1000) / 8
		fmt.Println("Tap CHange", lenghtDif)
		ResizeAndMove(fileName, 0, 0, 298+lenghtDif, 40)

		ResizeAndMove(size, 301+lenghtDif, 0, 62+lenghtDif, 40)

		ResizeAndMove(proBar, 365+(2*lenghtDif), 0, 145+(lenghtDif), 40)

		ResizeAndMove(seeders, 511+(3*lenghtDif), 0, 80+(lenghtDif), 40)

		ResizeAndMove(leechers, 592+(4*lenghtDif), 0, 80+(lenghtDif), 40)

		ResizeAndMove(downloadSpeed, 672+(5*lenghtDif), 0, 160+(lenghtDif), 40)

		ResizeAndMove(ETA, 832+(6*lenghtDif), 0, 80+(lenghtDif), 40)

		ResizeAndMove(dateOfAdd, 912+(7*lenghtDif), 0, 80+(lenghtDif), 40)

		win.Canvas().Content().Refresh()
		log.Println("Refresh:", win.Canvas().Size().Width)

	})

	list := widget.NewList(
		// lets change item count from 3 to 30
		func() int { return 1 }, // my list contain 3 items

		func() fyne.CanvasObject {

			//pro.Resize(fyne.NewSize(200, 50))

			//item1 := widget.NewAccordionItem("Attack On Titan S04Ep28", widget.NewLabel("The Episode"))

			//row := container.NewHBox(ac, canvas.NewRectangle(color.Black), size, pro, widget.NewLabel("11"), widget.NewLabel("18"), widget.NewLabel("5.4 MB/s"), widget.NewLabel("9 Min "), widget.NewLabel("6/2/2022"))
			//	container.NewAdaptiveGrid()

			row2 := container.NewWithoutLayout(fileName, size, proBar, seeders, leechers, downloadSpeed, ETA, dateOfAdd)

			return row2

		},
		// last one
		func(lii widget.ListItemID, co fyne.CanvasObject) {
			// update data of widget
			//co.(*widget.Label).SetText("Here is my Newtext")

			//co.Resize(fyne.NewSize(150, 50))
		},
	)
	return list, *item3
}

/*

&canvas.Line{
			Position1:   fyne.NewPos(130, 0),
			Position2:   fyne.NewPos(130, content.),
			StrokeWidth: 2,
			StrokeColor: color.White,
		},

}*/

func bottomInfo(win fyne.Window) fyne.CanvasObject {

	label1 := canvas.NewText("Linux Mint 21 Cinnamon 64Bit ISO", color.NRGBA{R: 50, G: 186, B: 219, A: 255})
	label1.TextSize = 14
	label1.TextStyle.Bold = true

	hashInfo := canvas.NewText("14bb29461c2243aa287cda681488e57c9eeb25af", color.NRGBA{R: 50, G: 186, B: 219, A: 255})
	hashInfo.TextSize = 14
	hashInfo.TextStyle.Bold = true

	downloaded := canvas.NewText("23863 downloads", color.NRGBA{R: 50, G: 186, B: 219, A: 255})
	downloaded.TextSize = 14
	downloaded.TextStyle.Bold = true

	pieceOFleft := canvas.NewText("7466/9851", color.NRGBA{R: 50, G: 186, B: 219, A: 255})
	pieceOFleft.TextSize = 14
	pieceOFleft.TextStyle.Bold = true

	fileList := canvas.NewText("Linux Mint 21 Cinnamon.exe", color.NRGBA{R: 50, G: 186, B: 219, A: 255})
	fileList.TextSize = 14
	fileList.TextStyle.Bold = true

	savedPath := canvas.NewText("/Downloads", color.NRGBA{R: 50, G: 186, B: 219, A: 255})
	savedPath.TextSize = 14
	savedPath.TextStyle.Bold = true

	info := container.NewGridWithColumns(2, container.NewHBox(widget.NewLabel("Torrent File:"), label1), container.NewHBox(widget.NewLabel("Hash:"), hashInfo), container.NewHBox(widget.NewLabel("Downloaded:"), downloaded), container.NewHBox(widget.NewLabel("Number of Pieces left:"), pieceOFleft), container.NewHBox(widget.NewLabel("File List:"), fileList), container.NewHBox(widget.NewLabel("Saved PATH:"), savedPath))
	peerList := widget.NewList(func() int { return 5 }, func() fyne.CanvasObject { return widget.NewLabel("198.565.256.87 Country: Antartica") }, func(lii widget.ListItemID, co fyne.CanvasObject) {})
	peerInfo := container.NewBorder(container.NewVBox(widget.NewLabel("Number of Connected Peers:"), widget.NewLabel("Peer List:")), nil, nil, nil, peerList)
	trackerList := widget.NewList(func() int { return 10 }, func() fyne.CanvasObject { return widget.NewLabel("udp://open.stealth.si:80/announce  ") }, func(lii widget.ListItemID, co fyne.CanvasObject) {})
	trackerInfo := container.NewBorder(container.NewVBox(widget.NewLabel("Number of Trackers:"), widget.NewLabel("Tracker List:")), nil, nil, nil, trackerList)

	trackerInfo.Hide()
	info.Hide()
	peerInfo.Hide()

	infoButton := widget.NewButton("Information", func() {

		if info.Visible() {
			info.Hide()
		} else {
			trackerInfo.Hide()
			peerInfo.Hide()
			info.Show()
		}
		win.Canvas().Refresh(info)
	})

	trackerButton := widget.NewButton("Tracker Stats", func() {
		if trackerInfo.Visible() {
			trackerInfo.Hide()
		} else {
			info.Hide()
			peerInfo.Hide()
			trackerInfo.Show()
		}
		win.Canvas().Refresh(trackerInfo)
	})

	peerInfoButton := widget.NewButton("Peers", func() {
		if peerInfo.Visible() {
			peerInfo.Hide()
		} else {

			trackerInfo.Hide()
			info.Hide()
			peerInfo.Show()
		}
		win.Canvas().Refresh(peerInfo)

	})
	//trackerList:=

	return container.NewBorder(container.NewVBox(container.NewHBox(infoButton, peerInfoButton, trackerButton)), nil, nil, nil, info, peerInfo, trackerInfo)

	//Information: TorrentName:, Size, how many GB left, Num of pieces left, Torrent hash, ALL the files in it list And Saved PATH
}
