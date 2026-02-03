// nolint:errcheck
package formats

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strings"

	"github.com/golang/freetype"
	"github.com/nfnt/resize"
	"go.ads.coffee/platform/server/internal/domain/ads"
	"go.ads.coffee/platform/server/internal/tools/filesystem"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
)

type Banner struct {
}

func NewBanner() *Banner {
	return &Banner{}
}

func (b *Banner) Banner(ctx context.Context, base string, banner ads.Banner, w http.ResponseWriter) error {
	image, err := filesystem.NewFileFromURL(ctx, banner.Image.Full(base))
	if err != nil {
		return err
	}

	buffer, err := image.Reader.Open()
	if err != nil {
		return err
	}

	defer buffer.Close()

	data, format, err := b.Render(buffer, banner.Description, banner.Title)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Length", fmt.Sprint(data.Len()))
	w.Header().Set("Content-Type", "image/"+format)

	if _, err = io.Copy(w, data); err != nil {
		return err
	}

	return nil
}

func (b *Banner) Render(file io.Reader, description, info string) (*bytes.Buffer, string, error) {
	// Параметры баннера (все размеры увеличены вдвое)
	width := 820                                       // было 410
	img := image.NewRGBA(image.Rect(0, 0, width, 360)) // было 180

	// Заливаем фон
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{248, 249, 250, 255}}, image.Point{}, draw.Src)

	// Загружаем шрифт
	f, err := freetype.ParseFont(goregular.TTF)
	if err != nil {
		panic(err)
	}

	// Контекст для текста (увеличиваем DPI для сохранения четкости)
	c := freetype.NewContext()
	c.SetDPI(72) // было 72
	c.SetFont(f)
	c.SetClip(img.Bounds())
	c.SetDst(img)
	c.SetSrc(image.Black)
	c.SetHinting(font.HintingNone)

	// Декодируем изображение
	data, format, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}

	// Масштабируем изображение до высоты 180 пикселей
	data = resize.Resize(0, 180, data, resize.Lanczos3)

	imgX, imgY := 30, 30 // было 15,15

	draw.Draw(img,
		image.Rect(imgX, imgY, imgX+data.Bounds().Dx(), imgY+data.Bounds().Dy()),
		data,
		image.Point{},
		draw.Src)

	// Разбиваем описание на строки
	lines := wrap(description, (width-imgX-data.Bounds().Dy()-30)/10)

	// Вычисляем общую высоту текста (межстрочный интервал увеличен вдвое)
	textHeight := len(lines) * 40 // было 20
	if textHeight < 1 {
		textHeight = 40
	}

	// Вычисляем стартовую позицию текста для центрирования
	textX := imgX + data.Bounds().Dx() + 30 // было +15
	textStartY := imgY + (data.Bounds().Dy()-textHeight)/2
	if textStartY < imgY {
		textStartY = imgY
	}

	// Рисуем основной текст с центрированием (размер шрифта увеличен вдвое)
	c.SetFontSize(30) // было 15
	for i, line := range lines {
		pt := freetype.Pt(textX, textStartY+40*i+30) // было 20*i+15
		c.DrawString(line, pt)
	}

	// Рисуем разделительную линию (толщина линии увеличена)
	markerY := imgY + data.Bounds().Dy() + 40                                           // было +20
	drawLine(img, 30, markerY-20, width-30, markerY-20, color.RGBA{220, 220, 220, 255}) // было 15,markerY-10

	// Информация о рекламодателе (размер шрифта увеличен вдвое)
	c.SetFontSize(20) // было 10
	c.SetSrc(image.NewUniform(color.RGBA{150, 150, 150, 255}))
	infoLines := wrap(info, 90)
	for _, line := range infoLines {
		markerY += 30                  // было 15
		pt := freetype.Pt(30, markerY) // было 15
		c.DrawString(line, pt)
	}

	// Сохраняем результат в том же формате, что и исходное изображение
	result := bytes.NewBuffer([]byte{})

	switch format {
	case "jpeg":
		err = jpeg.Encode(result, img, &jpeg.Options{Quality: 90}) // Качество JPEG
	case "png":
		err = png.Encode(result, img)
	// case "webp":
	// 	err = webp.Encode(outFile, rgba, &webp.Options{Lossless: true}) // Без потерь для WebP
	default:
		return nil, "", fmt.Errorf("undefined format: %s", format)
	}

	if err != nil {
		return nil, "", err
	}

	return result, format, nil
}

// Функция для разбивки текста на строки
func wrap(text string, lineLength int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= lineLength {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	lines = append(lines, currentLine)

	return lines
}

// Функция для рисования линии
func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	dx, dy := x2-x1, y2-y1
	steps := max(abs(dx), abs(dy))

	for i := 0; i <= steps; i++ {
		x := x1 + (dx * i / steps)
		y := y1 + (dy * i / steps)
		img.Set(x, y, col)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
