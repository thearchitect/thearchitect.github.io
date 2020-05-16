package ui

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/button"
	"github.com/gcla/gowid/widgets/columns"
	"github.com/gcla/gowid/widgets/edit"
	"github.com/gcla/gowid/widgets/fill"
	"github.com/gcla/gowid/widgets/framed"
	"github.com/gcla/gowid/widgets/holder"
	"github.com/gcla/gowid/widgets/pile"
	"github.com/gcla/gowid/widgets/shadow"
	"github.com/gcla/gowid/widgets/styled"
	"github.com/gcla/gowid/widgets/terminal"
	"github.com/gcla/gowid/widgets/text"
	"github.com/gdamore/tcell"
)

// https://github.com/gcla/gowid/blob/master/docs/Tutorial.md
// https://github.com/gcla/gowid/blob/master/docs/Widgets.md
// https://github.com/gcla/gowid/blob/master/docs/FAQ.md

func Run(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	palette := gowid.Palette{
		"invred": gowid.MakePaletteEntry(gowid.ColorBlack, gowid.ColorRed),
		"line":   gowid.MakeStyledPaletteEntry(gowid.NewUrwidColor("black"), gowid.NewUrwidColor("light gray"), gowid.StyleBold),
	}

	title := text.New(" Hello ")

	setTitle := func(t string, app gowid.IApp) {
		title.SetText(fmt.Sprintf(" %s ", t), app)
	}

	widgetHello := buildHello()
	widgetCrashRoom := buildCrashRoom(setTitle)
	widgetEmptiness := buildEmptiness()

	content := holder.New(widgetHello)
	description := text.New("", text.Options{Wrap: text.WrapAny})

	widgetMenu := buildMenu(
		buildButton("Hello", func(app gowid.IApp) {
			setTitle("Hello", app)
			description.SetText("", app)
			content.SetSubWidget(widgetHello, app)
		}),
		buildButton("Crash Room", func(app gowid.IApp) {
			setTitle("Crash Room", app)
			description.SetText(`Crash room (aka anger or rage room) is a place where people can vent their rage by destroying objects within a room. Here you can see the one and only Linux server crash room. It's free of charge and it's totally yours!`, app)
			content.SetSubWidget(widgetCrashRoom, app)
		}),
		buildButton("", func(app gowid.IApp) {
			// setTitle("", app)
			title.SetText("", app)
			description.SetText(`emptiness`, app)
			content.SetSubWidget(widgetEmptiness, app)
		}),
	)

	// vline := styled.New(fill.New('│'), gowid.MakePaletteRef("line"))
	// hline := styled.New(fill.New('⎯'), gowid.MakePaletteRef("line"))

	rightPanel := pile.New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{IWidget: widgetMenu, D: gowid.RenderFlow{}},
		// &gowid.ContainerWidget{IWidget: hline, D: gowid.RenderWithWeight{W: 1}},
		&gowid.ContainerWidget{IWidget: description, D: gowid.RenderWithWeight{W: 1}},
	})

	view := framed.New(columns.New([]gowid.IContainerWidget{
		&gowid.ContainerWidget{
			IWidget: framed.New(framed.New(content, framed.Options{
				Frame: framed.SpaceFrame,
			}), framed.Options{
				Frame:       framed.UnicodeAlt2Frame,
				TitleWidget: title,
			}),
			D: gowid.RenderWithWeight{W: 4},
		},
		&gowid.ContainerWidget{
			IWidget: fill.New(' '),
			D:       gowid.RenderWithUnits{U: 1},
		},
		&gowid.ContainerWidget{
			IWidget: framed.New(framed.New(rightPanel, framed.Options{
				Frame: framed.SpaceFrame,
			}), framed.Options{
				Frame: framed.UnicodeFrame,
			}),
			D: gowid.RenderWithWeight{W: 1},
		},
	}), framed.Options{
		Frame: framed.SpaceFrame,
	})

	app, err := gowid.NewApp(gowid.AppArgs{
		View:    view,
		Palette: &palette,
	})
	if err != nil {
		panic(err)
	}

	go func() {
		<-ctx.Done()
		app.Quit()
	}()

	app.SetColorMode(gowid.Mode256Colors)

	runApp(app, func(w, h int) {
		//log.Println("resize!", w, h)
	}, gowid.UnhandledInputFunc(func(app gowid.IApp, ev interface{}) bool {
		if evk, ok := ev.(*tcell.EventKey); ok {
			switch evk.Key() {
			case tcell.KeyCtrlC, tcell.KeyEsc, tcell.KeyCtrlQ:
				cancel()
				return true
			}
		}
		return false
	}))
}

func runApp(app *gowid.App, resize func(w, h int), unhandled gowid.IUnhandledInput) {
	defer app.Close()
	st := app.Runner()
	st.Start()
	defer st.Stop()
L:
	for {
		select {
		case ev := <-app.TCellEvents:
			app.HandleTCellEvent(ev, unhandled)

			switch ev := ev.(type) {
			case *tcell.EventResize:
				resize(ev.Size())
			}
		case ev := <-app.AfterRenderEvents:
			if ev == nil {
				break L
			}
			app.RunThenRenderEvent(ev)
		}
	}
}

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
				stld := styled.New(
					shdw,
					gowid.MakePaletteEntry(gowid.ColorGreen, gowid.ColorBlack),
				)
				return stld
			} else {
				txt := text.New(txt, text.Options{})
				frmd := framed.New(txt, framed.Options{Frame: framed.UnicodeFrame})
				shdw := shadow.New(frmd, 1)
				stld := styled.New(
					shdw,
					gowid.MakePaletteEntry(gowid.ColorRed, gowid.ColorBlack),
				)
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

func buildHello2() (gowid.IWidget, func(app gowid.IApp)) {
	w := edit.New(edit.Options{})
	//w := gowid.NewCanvas()
	//sh -c "figlet -f poison Hello | lolcat; while true; do sleep 2; done"
	return w, func(app gowid.IApp) {
		// poison sblood fraktur
		cmd := exec.Command("sh", "-c", "figlet -f poison Hello | lolcat")
		cmd.Stdout = &edit.Writer{Widget: w, IApp: app}
		//cmd.Stdout = w
		if err := cmd.Start(); err != nil {
			panic(err)
		}
		app.Redraw()
	}
}

func buildEmptiness() *fill.Widget {
	return fill.New(' ')
}

func buildHello() gowid.IWidget {
	// poison sblood fraktur
	termWidget, err := terminal.NewExt(terminal.Options{
		//Command:    []string{"/bin/sh", "-c", "figlet -f poison Hello | lolcat -t -a"},
		//Command:    []string{"figlet", "-f", "poison", "Hello"},
		Command: []string{"toilet", "-d", "/usr/local/share/figlet/fonts", "-f", "fraktur", "--filter", "metal", "Hello"},
		//Env: []string{"TERM=xterm-256color", "COLORTERM=truecolor"},
		Scrollback: 1024,
	})
	if err != nil {
		panic(err)
	}

	return termWidget
}

func buildCrashRoom(setTitle func(t string, app gowid.IApp)) *terminal.Widget {
	//tw := text.New(" Terminal ")
	//twi := styled.New(tw, gowid.MakePaletteRef("invred"))
	//twp := holder.New(tw)

	termWidget, err := terminal.NewExt(terminal.Options{
		Command: []string{"/bin/sh", "-l", "-c", "mc / /home/johnny; zsh -l"},
		//Env:               environment.Environ().Slice(),
		HotKeyPersistence: &terminal.HotKeyDuration{D: time.Second * 3},
		Scrollback:        8 * 1024,
	})
	if err != nil {
		panic(err)
	}

	{
		termWidget.OnProcessExited(gowid.WidgetCallback{Name: "cb",
			WidgetChangedFunction: func(app gowid.IApp, w gowid.IWidget) {
				app.Quit()
			},
		})
		termWidget.OnBell(gowid.WidgetCallback{Name: "cb",
			WidgetChangedFunction: func(app gowid.IApp, w gowid.IWidget) {
				go func() {
					time.Sleep(time.Millisecond * 800)
					app.Run(gowid.RunFunction(func(app gowid.IApp) {
						//title.SetText("", app)
					}))
				}()
			},
		})
		termWidget.OnSetTitle(gowid.WidgetCallback{Name: "cb",
			WidgetChangedFunction: func(app gowid.IApp, w gowid.IWidget) {
				setTitle((w.(*terminal.Widget)).GetTitle(), app)
			},
		})
	}

	return termWidget
}
