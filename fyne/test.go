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
Fix Header Bar


Make the first arrow and allign it, with color
make a function for the tree
make info,Tracker into buttons

make resizing the app adjusts everyhting
*/

func main() { // rage in the darkness
	// we can add styles later

	var filePath string

	myApp := app.New()
	myWindow := myApp.NewWindow("ReTorrent")
	myWindow.Resize(fyne.NewSize(1000, 500))

	appIcon, _ := fyne.LoadResourceFromPath("C:/Users/dontw/Downloads/nk257hl881b81.jpg")
	myApp.SetIcon(appIcon)

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

	item3 := fyne.NewMenuItem("Magent Link", nil)

	item2 := fyne.NewMenuItem("About", func() {

		link, err := url.Parse("https://www.youtube.com/watch?v=2JpkMXinO1M&list=RDSZOr9K01Eno")

		if err == nil {
			myApp.OpenURL(link)

		}

	})

	menu := fyne.NewMenu("File", item1, item3)

	menu2 := fyne.NewMenu("About", item2)

	myWindow.SetMainMenu(fyne.NewMainMenu(menu, menu2))

	// menu Done
	// 	split := container.NewVSplit(container.NewGridWithRows(2, bodyContainer(), container.NewMax(canvas.NewRectangle(color.Black), torrentList())), bottomInfo())
	seperator1 := widget.NewSeparator()
	seperator1.Resize(fyne.NewSize(3, 0))
	// listTets := container.NewMax(canvas.NewRectangle(color.Black), torrentList())

	// layoutTest := container.NewWithoutLayout(container.NewVBox(bodyContainer(myWindow.Canvas().Size()), seperator1), listTets)
	// listTets.Resize(fyne.NewSize(1000, 500))
	// listTets.Move(fyne.NewPos(0, 50))

	split := container.NewVSplit(container.NewGridWithRows(2, container.NewVBox(bodyContainer(myWindow.Canvas().Size()), seperator1), container.NewMax(canvas.NewRectangle(color.Black), torrentList())), bottomInfo())
	//split := container.NewVSplit(layoutTest, bottomInfo())

	split.Offset = 1.0
	myWindow.SetContent(split)

	log.Println("Hreres", split.Trailing.Position())

	myWindow.ShowAndRun()
}

func bodyContainer(size1 fyne.Size) fyne.CanvasObject {

	log.Println("SIZE:", size1)

	content := container.NewMax()
	size := &widget.Button{
		Alignment:  widget.ButtonAlignCenter,
		Text:       " Size ",
		OnTapped:   func() { fmt.Println("Tap Size") },
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
	listss := torrentList()
	listss.Resize(fyne.NewSize(400, 50))

	header := container.NewWithoutLayout(
		fileName,
		seperator,
		size, seperator2, progressBar, seperator3, Seeders, seperator4, Leechers, seperator5, downloadSpeed, seperator6, ETA, seperator7, date, seperator8,
		content)

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

	log.Println(header.Size())

	return header

}
func ResizeAndMove(element fyne.CanvasObject, xPos, yPos, width, height float32) {
	element.Resize(fyne.NewSize(width, height))
	element.Move(fyne.NewPos(xPos, yPos))
}

func torrentList() fyne.CanvasObject {

	sty := &fyne.TextStyle{
		Bold: true,
	}

	item := widget.NewAccordionItem("Attack On Titan S04Ep28", widget.NewLabel("The Episode"))
	fileName := widget.NewAccordion(item)
	size := widget.NewLabelWithStyle("1.07GB", fyne.TextAlignCenter, *sty)
	proBar := widget.NewProgressBar()
	seeders := widget.NewLabelWithStyle("8", fyne.TextAlignCenter, *sty)
	leechers := widget.NewLabelWithStyle("18", fyne.TextAlignCenter, *sty)
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
	proBar.SetValue(.92)

	list := widget.NewList(
		// lets change item count from 3 to 30
		func() int { return 30 }, // my list contain 3 items

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
	return list
}

/*

&canvas.Line{
			Position1:   fyne.NewPos(130, 0),
			Position2:   fyne.NewPos(130, content.),
			StrokeWidth: 2,
			StrokeColor: color.White,
		},


func headerRow() fyne.CanvasObject {
	col1Width := float32(25)
	col1X := theme.Padding()

	col2Width := (size.Width - col1Width) * 0.42
	col2X := col1X + col1Width + theme.Padding()

	col3Width := (size.Width - col2Width) * 0.25
	col3X := col2X + col2Width + theme.Padding()

	col4Width := (size.Width - col3Width) * 0.10
	col4X := col3X + col3Width + theme.Padding()

	col5Width := (size.Width - col4Width) * 0.13
	col5X := col4X + col4Width + theme.Padding()

	col6Width := (size.Width - col5X - col5Width)
	col6X := col5X + col5Width - theme.Padding()

	ResizeAndMove(objects[0], col1Width, col1X, l.maxMinSizeHeight)

}*/

func bottomInfo() fyne.CanvasObject {

	tree := widget.NewButton("Information", nil)
	tree2 := widget.NewButton("Peers", nil)
	tree3 := widget.NewButton("Tracker Stats", nil)

	return container.NewBorder(container.NewHBox(tree, tree2, tree3), nil, nil, nil)
}
