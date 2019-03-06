package salas

import (
	"fmt"
	"io"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func generateImage(salas *salasResponse, disabled int, stream *io.PipeWriter) {
	w := 700
	h := 500
	salah := float64(h / 5.0)
	totalHours := 22
	hourw := float64(w) / float64(totalHours)

	margin := 50.0

	c := gg.NewContext(w+int(2*margin), h+int(2*margin))

	c.SetRGB(1, 1, 1)
	c.Clear()

	// Draw everything as available
	c.SetRGB(19.0/255.0, 142.0/255.0, 0)
	c.DrawRectangle(margin, margin, float64(w), float64(h))
	c.Fill()

	if disabled > 22 {
		disabled = 0
	}

	if disabled > 0 {
		c.SetRGB(15.0/255.0, 53.0/255.0, 18.0/255.0)
		c.DrawRectangle(margin, margin, hourw*float64(disabled), float64(h))
		c.Fill()
	}

	font, err := truetype.Parse(goregular.TTF)

	if err != nil {
		panic(err)
	}

	face := truetype.NewFace(font, &truetype.Options{
		Size: 20,
	})

	subface := truetype.NewFace(font, &truetype.Options{
		Size: 15,
	})

	c.SetFontFace(face)
	for i, sala := range salas.Salas {
		for _, t := range sala.Occupied {
			si := timeToIndex(t.Start/100, t.Start%100)
			ei := timeToIndex(t.End/100, t.End%100)

			if si < disabled {
				c.SetRGB(68.0/255.0, 19.0/255.0, 19.0/255.0)
			} else {
				c.SetRGB(142.0/255.0, 0, 0)
			}
			c.DrawRectangle(hourw*float64(si)+margin, salah*float64(i)+margin, float64(float64(ei-si)*hourw), salah)
			c.Fill()
		}

		c.SetRGB(1, 1, 1)

		for j := 0; j < totalHours; j++ {
			c.DrawRectangle(float64(j)*hourw+margin, salah*float64(i)+margin, hourw, salah)
			c.SetLineWidth(1.0)
			c.Stroke()
		}

		// Draw sala text
		c.SetRGB(0, 0, 0)
		c.RotateAbout(gg.Radians(-90), margin, float64(salah*float64(i)+salah/2.0)+margin)
		c.DrawStringAnchored(fmt.Sprintf("Sala #%v", sala.ID), margin, float64(salah*float64(i)+salah/3.0)+margin, 0.5, 0.5)
		c.RotateAbout(gg.Radians(90), margin, float64(salah*float64(i)+salah/2.0)+margin)
	}

	// Draw hour labels
	c.SetRGB(0, 0, 0)
	c.SetFontFace(subface)
	for j := 0; j < totalHours; j++ {
		c.RotateAbout(gg.Radians(-45), float64(j)*hourw+2+margin, 45)
		c.DrawString(indexToTime(j), float64(j)*hourw+2+margin, 45)
		c.RotateAbout(gg.Radians(45), float64(j)*hourw+2+margin, 45)
	}

	c.SetFontFace(face)

	// Draw legend
	legspacing := 87.0
	legw := 50.0

	c.SetRGB(19.0/255.0, 142.0/255.0, 0)
	c.DrawRectangle(margin, margin+float64(h)+10, legw, 30)
	c.Fill()
	c.SetRGB(0, 0, 0)
	c.DrawString("Libre", margin+legw+5, margin+float64(h)+margin-17)

	c.SetRGB(142.0/255.0, 0, 0)
	c.DrawRectangle(margin+legspacing*2, margin+float64(h)+10, legw, 30)
	c.Fill()
	c.SetRGB(0, 0, 0)
	c.DrawString("Reservado", margin+legspacing*2+legw+5, margin+float64(h)+margin-17)

	c.SetRGB(68.0/255.0, 19.0/255.0, 19.0/255.0)
	c.DrawRectangle(margin+legspacing*4, margin+float64(h)+10, legw, 30)
	c.Fill()
	c.SetRGB(0, 0, 0)
	c.DrawString("Ocupado", margin+legspacing*4+legw+5, margin+float64(h)+margin-17)

	c.SetRGB(15.0/255.0, 53.0/255.0, 18.0/255.0)
	c.DrawRectangle(margin+legspacing*6, margin+float64(h)+10, legw, 30)
	c.Fill()
	c.SetRGB(0, 0, 0)
	c.DrawString("Pasado", margin+legspacing*6+legw+5, margin+float64(h)+margin-17)

	c.EncodePNG(stream)
	stream.Close()
}
