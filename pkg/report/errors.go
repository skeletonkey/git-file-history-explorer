package report

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func PanicOnError(e error) {
	if e != nil {
		panic(e)
	}
}

func ErrorPopUp(err error, w fyne.Window, closeFn func()) {
	errDialog := dialog.NewError(err, w)
	errDialog.SetDismissText("Acknowlege")
	if closeFn != nil {
		errDialog.SetOnClosed(closeFn)
	}
	errDialog.Show()
}
