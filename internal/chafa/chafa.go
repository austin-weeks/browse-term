// Package chafa wraps the chafa library.
package chafa

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"

	"github.com/ploMP4/chafa-go"
)

// const symbols = chafa.CHAFA_SYMBOL_TAG_SEXTANT
const symbols = chafa.CHAFA_SYMBOL_TAG_SEXTANT

var config *chafa.CanvasConfig

// Code adapted from chafa-go docs
func init() {
	// Specify the symbols we want
	symbolMap := chafa.SymbolMapNew()
	chafa.SymbolMapAddByTags(symbolMap, symbols)
	// Set up a configuration with the symbols
	config = chafa.CanvasConfigNew()
	chafa.CanvasConfigSetSymbolMap(config, symbolMap)
}

func setSize(w, h int) {
	chafa.CanvasConfigSetGeometry(config, int32(w), int32(h))
}

func RenderImage(data []byte, w int) (res string, err error) {
	// TODO: aspect ratio isn't quite right
	w = 10
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("failed to render image with chafa: %v", r)
		}
	}()

	// Decode image
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Convert to RGBA
	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	// Set canvas size based on aspect ratio
	imgW, imgH := bounds.Dx(), bounds.Dy()
	canvasH := calculateCanvasHeight(imgW, imgH, w)
	setSize(w, canvasH)

	canvas := chafa.CanvasNew(config)
	defer chafa.CanvasUnRef(canvas)

	// Draw pixels
	chafa.CanvasDrawAllPixels(
		canvas,
		chafa.CHAFA_PIXEL_RGBA8_UNASSOCIATED,
		rgba.Pix, // Use rgba.Pix
		int32(imgW),
		int32(imgH),
		int32(imgW*4), // stride: width * 4 for RGBA
	)

	// Generate a string that will show the canvas contents on a terminal
	gs := chafa.CanvasPrint(canvas, nil)

	// TODO: add some padding on the left for the images
	return gs.String(), nil
}

func calculateCanvasHeight(imgW, imgH, outputW int) int {
	h := int(float64(imgH) * float64(outputW) / float64(imgW) / 2.0)
	if h == 0 {
		h = 1
	}
	return h
}
