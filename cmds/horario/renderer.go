package horario

import (
	"fmt"
	"io"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func generateImage(horario [][]class, stream *io.PipeWriter) {
	margin := 50.0

	w := 1000
	h := 700
	colW := float64(w) / float64(6)

	c := gg.NewContext(w+int(2*margin), h+int(2*margin))

	c.SetRGB(1, 1, 1)
	c.Clear()

	// Get fonts
	font, err := truetype.Parse(goregular.TTF)

	if err != nil {
		panic(err)
	}

	tableHeaders := []string{"HORA", "LUNES", "MARTES", "MIÉRCOLES", "JUEVES", "VIERNES"}
	tableData := [][]string{}
	hasRow := []bool{}

	// Generate table
	for i := 9; i <= 20; i++ {
		tableData = append(tableData, []string{fmt.Sprintf("%v:00\n%v:00", i, i+1), "", "", "", "", ""})
		hasRow = append(hasRow, false)
	}

	// Populate table
	for day := 0; day < len(horario); day++ {
		for _, clase := range horario[day] {
			for i := clase.Start; i < clase.End; i++ {
				tableData[i-9][day+1] = clase.Name
				hasRow[i-9] = true
			}
		}
	}

	// Expunge all rows without data
	for i := len(hasRow) - 1; i >= 0; i-- {
		if !hasRow[i] {
			copy(tableData[i:], tableData[i+1:])
			tableData[len(tableData)-1] = nil
			tableData = tableData[:len(tableData)-1]
		}
	}

	// Draw the table

	face := truetype.NewFace(font, &truetype.Options{
		Size: 35,
	})

	subface := truetype.NewFace(font, &truetype.Options{
		Size: 26,
	})

	rowH := float64(h) / float64(len(tableData)+1)

	// Draw headers
	c.SetFontFace(subface)
	for i := range tableHeaders {
		x := colW*float64(i) + margin
		y := margin

		c.DrawRectangle(x, y, colW, rowH)
		c.SetRGB(192.0/255.0, 214.0/255.0, 235.0/255.0)
		c.Fill()

		c.SetRGB(0.0, 0.0, 0.0)
		c.DrawRectangle(x, y, colW, rowH)
		c.SetLineWidth(2.0)
		c.Stroke()

		c.SetRGB(0.0, 0.0, 0.0)
		c.DrawStringAnchored(tableHeaders[i], x+colW/2.0, y+rowH/2.0, 0.5, 0.5)
	}

	c.SetFontFace(face)
	for i := range tableData {
		y := margin + rowH*float64(i+1)

		for j := range tableData[i] {
			x := margin + colW*float64(j)

			// Text

			// Draw time on 2 lines
			if j == 0 {
				c.DrawRectangle(x, y, colW, rowH)
				c.SetRGB(192.0/255.0, 214.0/255.0, 235.0/255.0)
				c.Fill()
				c.SetRGB(0.0, 0.0, 0.0)
				parts := strings.Split(tableData[i][j], "\n")

				_, strh := c.MeasureString(tableData[i][j])

				c.DrawStringAnchored(parts[0], x+colW/2.0, y+rowH/2.0-strh/2.0, 0.5, 0.5)
				c.DrawStringAnchored(parts[1], x+colW/2.0, y+rowH/2.0+strh/2.0, 0.5, 0.5)

			} else {
				c.DrawStringAnchored(tableData[i][j], x+colW/2.0, y+rowH/2.0, 0.5, 0.5)
			}

			// Outer rect
			c.SetRGB(0.0, 0.0, 0.0)
			c.DrawRectangle(x, y, colW, rowH)
			c.SetLineWidth(2.0)
			c.Stroke()
		}
	}

	c.EncodePNG(stream)
	stream.Close()
}
