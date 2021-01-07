package main

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type iconAndTextButton struct {
	theme *material.Theme
}
type imageAndTextButton struct {
	theme *material.Theme
}
type imageAndTextAndTagsButton struct {
	theme *material.Theme
}
type textButton struct {
	theme *material.Theme
}

func (b textButton) Layout(gtx layout.Context, button *widget.Clickable, color color.NRGBA, word string) D {
	l := material.ButtonLayout(b.theme, button)
	l.Background = color
	return l.Layout(gtx, func(gtx C) D {

		return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {

			b.theme.Palette.Bg = color
			iconAndLabel := layout.Flex{Axis: layout.Horizontal}

			layLabel := layout.Rigid(func(gtx C) D {
				return widget.Label{}.Layout(gtx, b.theme.Shaper, text.Font{}, b.theme.TextSize, word)
			})

			return iconAndLabel.Layout(gtx, layLabel)

		})
	})
}

func (b iconAndTextButton) Layout(gtx layout.Context, button *widget.Clickable, icon *widget.Icon, color color.NRGBA, word string) D {
	l := material.ButtonLayout(b.theme, button)
	l.Background = color

	return l.Layout(gtx, func(gtx C) D {

		return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {

			iconAndLabel := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
			var textIconSpacer unit.Value
			if word == "" {
				textIconSpacer = unit.Dp(0)
			} else {
				textIconSpacer = unit.Dp(5)
			}

			layIcon := layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: textIconSpacer}.Layout(gtx, func(gtx C) D {
					if icon != nil {
						return icon.Layout(gtx, unit.Dp(20))
					}
					return D{}
				})
			})

			layLabel := layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: textIconSpacer}.Layout(gtx, func(gtx C) D {
					return widget.Label{}.Layout(gtx, b.theme.Shaper, text.Font{}, b.theme.TextSize, word)
				})
			})

			return iconAndLabel.Layout(gtx, layIcon, layLabel)

		})
	})
}

func (b imageAndTextButton) Layout(gtx layout.Context, button *widget.Clickable, image *paint.ImageOp, color color.NRGBA, word string) D {
	l := material.ButtonLayout(b.theme, button)
	l.Background = color

	return l.Layout(gtx, func(gtx C) D {

		return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {

			iconAndLabel := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
			var textIconSpacer unit.Value
			if word == "" {
				textIconSpacer = unit.Dp(0)
			} else {
				textIconSpacer = unit.Dp(5)
			}

			layIcon := layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: textIconSpacer}.Layout(gtx, func(gtx C) D {
					if image != nil {
						return widget.Image{Src: *image, Scale: 1}.Layout(gtx)
					}
					return D{}
				})
			})

			layLabel := layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: textIconSpacer}.Layout(gtx, func(gtx C) D {
					return widget.Label{}.Layout(gtx, b.theme.Shaper, text.Font{}, b.theme.TextSize, word)
				})
			})

			return iconAndLabel.Layout(gtx, layIcon, layLabel)

		})
	})
}

func (b imageAndTextAndTagsButton) Layout(gtx layout.Context, button *widget.Clickable, image *paint.ImageOp, color color.NRGBA, word string, tags []string) D {
	l := material.ButtonLayout(b.theme, button)
	l.Background = color

	return l.Layout(gtx, func(gtx C) D {

		return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx C) D {

			iconAndLabel := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
			var textIconSpacer unit.Value
			if word == "" {
				textIconSpacer = unit.Dp(0)
			} else {
				textIconSpacer = unit.Dp(5)
			}

			layTags := layout.Rigid(func(gtx C) D {
				l := layout.List{Axis: layout.Vertical}
				return l.Layout(gtx, len(tags), func(gtx layout.Context, j int) D {
					th := material.Label(gui.Theme, unit.Dp(12), tags[j])
					th.Font.Style = text.Italic
					return th.Layout(gtx)
				})
			})

			layIcon := layout.Rigid(func(gtx C) D {
				return layout.Inset{Right: textIconSpacer}.Layout(gtx, func(gtx C) D {
					if image != nil {
						return widget.Image{Src: *image, Scale: 1}.Layout(gtx)
					}
					return D{}
				})
			})

			layLabel := layout.Rigid(func(gtx C) D {
				return layout.Inset{Left: textIconSpacer}.Layout(gtx, func(gtx C) D {
					return widget.Label{}.Layout(gtx, b.theme.Shaper, text.Font{}, b.theme.TextSize, word)
				})
			})

			return iconAndLabel.Layout(gtx, layIcon, layLabel, layTags)

		})
	})
}
