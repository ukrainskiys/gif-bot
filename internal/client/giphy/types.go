package giphy

type GifType string

const STICKER = GifType("sticker")
const GIF = GifType("gif")

func (gt GifType) toPath() string {
	return string(gt + "s")
}
