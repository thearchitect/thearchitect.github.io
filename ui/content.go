package ui

import (
	"github.com/gcla/gowid"
	"github.com/gcla/gowid/widgets/text"
	"github.com/thearchitect/go-figure"
)

type Content []text.ContentSegment

func (content Content) Banner(txt, font string, palette []gowid.ICellStyler) Content {
	if len(font) == 0 {
		font = "sblood"
	}

	if len(palette) == 0 {
		palette = []gowid.ICellStyler{
			gowid.MakePaletteEntry(gowid.ColorLightGreen, gowid.ColorNone),
			gowid.MakePaletteEntry(gowid.ColorGreen, gowid.ColorNone),
			gowid.MakePaletteEntry(gowid.ColorDarkGreen, gowid.ColorNone),
			gowid.MakePaletteEntry(gowid.ColorLightGray, gowid.ColorNone),
		}
	}

	banner := figure.NewFigure(txt, font, false).Slicify()

	for y, row := range banner {
		for x, col := range row {
			content = append(content, text.ContentSegment{
				Text:  string(col),
				Style: palette[(x+y)%len(palette)],
			})
		}
		content = append(content, text.ContentSegment{Text: "\n"})
	}

	return content
}

func (content Content) Text(txt string, style gowid.ICellStyler) Content {
	if style == nil {
		style = gowid.MakePaletteEntry(gowid.ColorGreen, gowid.ColorNone)
	}

	txt += "\n"

	content = append(content, text.ContentSegment{
		Text:  txt,
		Style: style,
	})
	return content
}

func (content Content) Content() *text.Content {
	return text.NewContent(content)
}
