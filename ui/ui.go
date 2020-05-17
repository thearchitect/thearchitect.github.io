package ui

import (
	"context"
	"fmt"
	"time"

	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/columns"
	"github.com/gcla/gowid/widgets/fill"
	"github.com/gcla/gowid/widgets/framed"
	"github.com/gcla/gowid/widgets/holder"
	"github.com/gcla/gowid/widgets/pile"
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
			description.SetText(`todo`, app)
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
	}, columns.Options{}), framed.Options{
		Frame: framed.SpaceFrame,
	})

	app, err := gowid.NewApp(gowid.AppArgs{
		View: view,
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

func buildHello() gowid.IWidget {
	c := (Content{}).
		Banner("HELLO", "sblood", nil).
		Text("", nil).
		Text("", nil).
		Text("1234567890", nil).
		Text("", nil).
		Text("", nil).
		Content()
	w := text.NewFromContentExt(c, text.Options{
		Wrap: text.WrapClip,
	})
	return w
}

func buildEmptiness() gowid.IWidget {
	//dim := gowid.RenderFixed{}
	//w := columns.New([]gowid.IContainerWidget{
	//	&gowid.ContainerWidget{
	//		IWidget: text.New("WAKE"),
	//		D:       dim,
	//	},
	//	&gowid.ContainerWidget{
	//		IWidget: text.New("UP"),
	//		D:       dim,
	//	},
	//	&gowid.ContainerWidget{
	//		IWidget: text.New("NEO"),
	//		D:       dim,
	//	},
	//	&gowid.ContainerWidget{
	//		IWidget: text.New("FOLLOW"),
	//		D:       dim,
	//	},
	//	&gowid.ContainerWidget{
	//		IWidget: text.New("THE"),
	//		D:       dim,
	//	},
	//	&gowid.ContainerWidget{
	//		IWidget: text.New("WHITE"),
	//		D:       dim,
	//	},
	//	&gowid.ContainerWidget{
	//		IWidget: text.New("RABBIT"),
	//		D:       dim,
	//	},
	//}, columns.Options{
	//	Wrap: true,
	//})

	return fill.New('•')
}

func buildCrashRoom(setTitle func(t string, app gowid.IApp)) *terminal.Widget {
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
