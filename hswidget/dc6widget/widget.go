package dc6widget

import (
	"fmt"
	"log"

	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
)

const (
	// nolint:gomnd // it was constant
	maxAlpha = uint8(255)
)

// widget represents dc6viewer's widget
type widget struct {
	id            string
	dc6           *d2dc6.DC6
	textureLoader hscommon.TextureLoader
	palette       *[256]d2interface.Color
}

// Create creates new widget
func Create(state []byte, palette *[256]d2interface.Color, textureLoader hscommon.TextureLoader, id string, dc6 *d2dc6.DC6) giu.Widget {
	result := &widget{
		id:            id,
		dc6:           dc6,
		textureLoader: textureLoader,
		palette:       palette,
	}

	if giu.Context.GetState(result.getStateID()) == nil && state != nil {
		s := result.getState()
		s.Decode(state)
		result.setState(s)
	}

	return result
}

// Build builds a widget
func (p *widget) Build() {
	state := p.getState()

	// nolint:gocritic // that's for now, will be more cases
	switch state.mode {
	case dc6WidgetViewer:
		p.makeViewerLayout().Build()
	}
}

func (p *widget) makeViewerLayout() giu.Layout {
	viewerState := p.getState()

	imageScale := uint32(viewerState.controls.scale)
	curFrameIndex := int(viewerState.controls.frame) + (int(viewerState.controls.direction) * int(p.dc6.FramesPerDirection))
	dirIdx := int(viewerState.controls.direction)

	textureIdx := dirIdx*int(p.dc6.FramesPerDirection) + int(viewerState.controls.frame)

	if imageScale < 1 {
		imageScale = 1
	}

	err := giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
	if err != nil {
		log.Print(err)
	}

	w := float32(p.dc6.Frames[curFrameIndex].Width * imageScale)
	h := float32(p.dc6.Frames[curFrameIndex].Height * imageScale)

	var widget *giu.ImageWidget
	if viewerState.textures == nil || len(viewerState.textures) <= int(viewerState.controls.frame) ||
		viewerState.textures[curFrameIndex] == nil {
		widget = giu.Image(nil).Size(w, h)
	} else {
		widget = giu.Image(viewerState.textures[textureIdx]).Size(w, h)
	}

	return giu.Layout{
		giu.Label(fmt.Sprintf(
			"Version: %v\t Flags: %b\t Encoding: %v\t",
			p.dc6.Version,
			int64(p.dc6.Flags),
			p.dc6.Encoding,
		)),
		giu.Label(fmt.Sprintf("Directions: %v\tFrames per Direction: %v", p.dc6.Directions, p.dc6.FramesPerDirection)),
		giu.Custom(func() {
			imgui.BeginGroup()
			if p.dc6.Directions > 1 {
				imgui.SliderInt("Direction", &viewerState.controls.direction, 0, int32(p.dc6.Directions-1))
			}

			if p.dc6.FramesPerDirection > 1 {
				imgui.SliderInt("Frames", &viewerState.controls.frame, 0, int32(p.dc6.FramesPerDirection-1))
			}

			imgui.SliderInt("Scale", &viewerState.controls.scale, 1, 8)

			imgui.EndGroup()
		}),
		giu.Separator(),
		widget,
	}
}
