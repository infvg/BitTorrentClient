/*package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Image")

	r, _ := fyne.LoadResourceFromURLString("https://m.media-amazon.com/images/M/MV5BMmQ3NzVmOWUtNTQ5Yi00MDczLWIxMDYtNGU5ZGU2YjI1NDY2XkEyXkFqcGdeQVRoaXJkUGFydHlJbmdlc3Rpb25Xb3JrZmxvdw@@._V1_.jpg")
	image := canvas.NewImageFromResource(r)

	//image := canvas.NewImageFromResource(theme.FyneLogo())
	//image := canvas.NewImageFromURI(link)
	// image := canvas.NewImageFromImage(src)
	// image := canvas.NewImageFromReader(reader, name)
	// image := canvas.NewImageFromFile(fileName)
	image.FillMode = canvas.ImageFillOriginal
	w.SetContent(image)

	w.ShowAndRun()
}
*/