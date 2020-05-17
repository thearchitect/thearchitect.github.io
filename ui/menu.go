package ui

import (
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/button"
	"github.com/gcla/gowid/widgets/framed"
	"github.com/gcla/gowid/widgets/holder"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/shadow"
	"github.com/gcla/gowid/widgets/styled"
	"github.com/gcla/gowid/widgets/text"
)

type Button struct {
	Name     string
	Callback func(app gowid.IApp)
}

func buildButton(name string, cb func(app gowid.IApp)) Button {
	return Button{
		Name:     name,
		Callback: cb,
	}
}

func buildMenu(buttons ...Button) *pile.Widget {
	var (
		activate   = make([]func(app gowid.IApp), len(buttons))
		deactivate = make([]func(app gowid.IApp), len(buttons))
	)

	buildButton := func(n int, txt string, cb func(app gowid.IApp)) gowid.IWidget {
		buildButton := func(active bool) gowid.IWidget {
			if active {
				txt := text.New(txt, text.Options{})
				frmd := framed.New(txt, framed.Options{Frame: framed.UnicodeFrame})
				shdw := shadow.New(frmd, 1)
				stld := styled.New(shdw, gowid.MakePaletteEntry(gowid.ColorGreen, gowid.ColorNone))
				return stld
			} else {
				txt := text.New(txt, text.Options{})
				frmd := framed.New(txt, framed.Options{Frame: framed.UnicodeFrame})
				shdw := shadow.New(frmd, 1)
				stld := styled.New(shdw, gowid.MakePaletteEntry(gowid.ColorRed, gowid.ColorNone))
				return stld
			}
		}

		activeButton := buildButton(true)
		inactiveButton := buildButton(false)

		var btnw gowid.IWidget
		if n == 0 {
			btnw = activeButton
		} else {
			btnw = inactiveButton
		}
		hldr := holder.New(btnw)
		btn := button.New(hldr, button.Options{})

		activate[n] = func(app gowid.IApp) {
			hldr.SetSubWidget(activeButton, app)
		}
		deactivate[n] = func(app gowid.IApp) {
			hldr.SetSubWidget(inactiveButton, app)
		}

		handleActivation := func(app gowid.IApp) {
			for i := 0; i < len(buttons); i++ {
				if i == n {
					activate[i](app)
				} else {
					deactivate[i](app)
				}
			}
		}

		btn.OnClick(gowid.WidgetCallback{
			WidgetChangedFunction: func(app gowid.IApp, w gowid.IWidget) {
				handleActivation(app)
				cb(app)
			},
		})

		return btn
	}

	var widgets []gowid.IContainerWidget

	for n, b := range buttons {
		widgets = append(widgets, &gowid.ContainerWidget{
			IWidget: buildButton(n, b.Name, b.Callback),
			D:       gowid.RenderFlow{},
		})
	}

	return pile.New(widgets, pile.Options{})
}
