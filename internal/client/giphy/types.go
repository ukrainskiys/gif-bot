package giphy

import "strings"

type GifType string

const STICKER = GifType("sticker")
const GIF = GifType("gif")

func (gt GifType) toPath() string {
	return string(gt + "s")
}

func (gt GifType) String() string {
	return strings.ToUpper(string(gt))
}

func ParseType(src string) GifType {
	if src == "gif" || src == GIF.String() {
		return GIF
	} else if src == "sticker" || src == STICKER.String() {
		return STICKER
	} else {
		return "null"
	}
}
