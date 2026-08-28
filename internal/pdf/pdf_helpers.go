package pdf

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	stdcolor "image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"strings"
	"sync"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"github.com/jung-kurt/gofpdf"
	"raven.go.invoice-builder/internal/models"
	"raven.go.invoice-builder/ui"
)

type paperConfig struct {
	Format       string
	IsNarrow     bool
	IsExtraSmall bool
	IsA5         bool
	TitleSize    float64
	BodySize     float64
	SmallSize    float64
	DescCol      uint
	QtyCol       uint
	PriceCol     uint
	AmountCol    uint
}

func getPaperConfig(paperSizeStr string) (pdf.Maroto, paperConfig) {
	format := strings.ToLower(strings.TrimSpace(paperSizeStr))
	switch format {
	case "pos_58", "pos58", "receipt_58", "58mm":
		m := pdf.NewMarotoCustomSize(consts.Portrait, consts.PageSize("58mm"), "mm", 58.0, 180.0)
		m.SetPageMargins(2, 3, 2)
		return m, paperConfig{
			Format:       "pos_58",
			IsNarrow:     true,
			IsExtraSmall: true,
			TitleSize:    10,
			BodySize:     8,
			SmallSize:    7,
			DescCol:      4,
			QtyCol:       2,
			PriceCol:     3,
			AmountCol:    3,
		}
	case "pos_80", "pos80", "receipt_80", "80mm", "bill":
		m := pdf.NewMarotoCustomSize(consts.Portrait, consts.PageSize("80mm"), "mm", 80.0, 220.0)
		m.SetPageMargins(4, 4, 4)
		return m, paperConfig{
			Format:       "pos_80",
			IsNarrow:     true,
			IsExtraSmall: false,
			TitleSize:    12,
			BodySize:     9,
			SmallSize:    8,
			DescCol:      4,
			QtyCol:       2,
			PriceCol:     3,
			AmountCol:    3,
		}
	case "a5":
		m := pdf.NewMaroto(consts.Portrait, consts.A5)
		m.SetPageMargins(8, 10, 8)
		return m, paperConfig{
			Format:       "a5",
			IsNarrow:     false,
			IsExtraSmall: false,
			IsA5:         true,
			TitleSize:    15,
			BodySize:     10,
			SmallSize:    9,
			DescCol:      5,
			QtyCol:       2,
			PriceCol:     2,
			AmountCol:    3,
		}
	case "letter":
		m := pdf.NewMaroto(consts.Portrait, consts.Letter)
		m.SetPageMargins(10, 15, 10)
		return m, paperConfig{
			Format:       "letter",
			IsNarrow:     false,
			TitleSize:    18,
			BodySize:     11,
			SmallSize:    9,
			DescCol:      5,
			QtyCol:       2,
			PriceCol:     2,
			AmountCol:    3,
		}
	case "legal":
		m := pdf.NewMaroto(consts.Portrait, consts.Legal)
		m.SetPageMargins(10, 15, 10)
		return m, paperConfig{
			Format:       "legal",
			IsNarrow:     false,
			TitleSize:    18,
			BodySize:     11,
			SmallSize:    9,
			DescCol:      5,
			QtyCol:       2,
			PriceCol:     2,
			AmountCol:    3,
		}
	default: // "a4" or fallback
		m := pdf.NewMaroto(consts.Portrait, consts.A4)
		m.SetPageMargins(10, 15, 10)
		return m, paperConfig{
			Format:       "a4",
			IsNarrow:     false,
			TitleSize:    18,
			BodySize:     11,
			SmallSize:    9,
			DescCol:      5,
			QtyCol:       2,
			PriceCol:     2,
			AmountCol:    3,
		}
	}
}

func getPrimaryColor() color.Color {
	return color.Color{
		Red:   74,
		Green: 124,
		Blue:  89,
	}
}

func getDarkColor() color.Color {
	return color.Color{
		Red:   46,
		Green: 50,
		Blue:  48,
	}
}

func getSecondaryColor() color.Color {
	return color.Color{
		Red:   116,
		Green: 121,
		Blue:  110,
	}
}

var (
	fadedBrandLogoCache *bytes.Buffer
	fadedBrandLogoOnce  sync.Once
)

func getFadedBrandLogoBuffer() *bytes.Buffer {
	fadedBrandLogoOnce.Do(func() {
		f, err := ui.Files.Open("asset/img/brand_logo.png")
		if err != nil {
			return
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			return
		}

		bounds := img.Bounds()
		mask := image.NewUniform(stdcolor.Alpha{8}) // ~3% opacity (255 * 0.03 = 7.6)
		faded := image.NewRGBA(bounds)
		draw.DrawMask(faded, bounds, img, image.Point{}, mask, image.Point{}, draw.Over)

		var buf bytes.Buffer
		err = png.Encode(&buf, faded)
		if err == nil {
			fadedBrandLogoCache = &buf
		}
	})
	return fadedBrandLogoCache
}

func buildHeader(m pdf.Maroto, cfg paperConfig, profile *models.BusinessProfile, logoPath, companyName, title string, wideLines []string, narrowLines []string) {
	// Load and fade the brand logo for the backdrop
	fadedLogoBuf := getFadedBrandLogoBuffer()
	if fadedLogoBuf != nil {
		if pdfM, ok := m.(*pdf.PdfMaroto); ok {
			opts := gofpdf.ImageOptions{
				ImageType:             "PNG",
				ReadDpi:               true,
				AllowNegativePosition: true,
			}
			pdfM.Pdf.RegisterImageOptionsReader("watermark", opts, fadedLogoBuf)
		}
	}

	m.RegisterHeader(func() {
		// Draw watermark backdrop first so it stays behind everything
		if fadedLogoBuf != nil {
			if pdfM, ok := m.(*pdf.PdfMaroto); ok {
				width, _ := m.GetPageSize()
				imageWidth := 100.0
				yPos := 100.0
				if cfg.IsNarrow {
					imageWidth = 50.0
					yPos = 50.0
				}
				xPos := (width - imageWidth) / 2
				pdfM.Pdf.ImageOptions("watermark", xPos, yPos, imageWidth, 0, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true, AllowNegativePosition: true}, 0, "")
			}
		}

		if logoPath != "" {
			if cfg.IsNarrow {
				m.Row(10.0, func() {
					m.Col(12, func() {
						_ = m.FileImage(logoPath, props.Rect{
							Percent: 50,
							Center:  true,
						})
					})
				})
			} else {
				m.Row(14.0, func() {
					m.Col(4, func() {
						_ = m.FileImage(logoPath, props.Rect{
							Percent: 100,
							Center:  false,
						})
					})
				})
			}
		}

		if cfg.IsNarrow {
			m.Row(10.0, func() {
				m.Col(12, func() {
					m.Text(companyName, props.Text{
						Size:  cfg.TitleSize,
						Style: consts.Bold,
						Align: consts.Center,
						Color: getDarkColor(),
					})
				})
			})
			if profile != nil && (profile.Address != "" || profile.RegistrationNumber != "") {
				m.Row(8.0, func() {
					m.Col(12, func() {
						info := profile.Address
						if profile.RegistrationNumber != "" {
							if info != "" {
								info += " | Reg #: " + profile.RegistrationNumber
							} else {
								info = "Reg #: " + profile.RegistrationNumber
							}
						}
						m.Text(info, props.Text{
							Size:  cfg.SmallSize,
							Align: consts.Center,
							Color: getSecondaryColor(),
						})
					})
				})
			}
			m.Row(10.0, func() {
				m.Col(12, func() {
					m.Text(title, props.Text{
						Size:  cfg.TitleSize - 1,
						Style: consts.Bold,
						Align: consts.Center,
						Color: getPrimaryColor(),
					})
				})
			})
			for _, line := range narrowLines {
				m.Row(8.0, func() {
					m.Col(12, func() {
						m.Text(line, props.Text{
							Size:  cfg.SmallSize,
							Align: consts.Center,
							Color: getSecondaryColor(),
						})
					})
				})
			}
		} else {
			m.Row(28.0, func() {
				m.Col(6, func() {
					m.Text(companyName, props.Text{
						Size:  cfg.TitleSize,
						Style: consts.Bold,
						Color: getDarkColor(),
					})
					if profile != nil && profile.Address != "" {
						m.Text(profile.Address, props.Text{
							Top:   9,
							Size:  cfg.BodySize,
							Color: getSecondaryColor(),
						})
					}
					if profile != nil && profile.RegistrationNumber != "" {
						m.Text(fmt.Sprintf("Reg #: %s", profile.RegistrationNumber), props.Text{
							Top:   15,
							Size:  cfg.SmallSize,
							Color: getSecondaryColor(),
						})
					}
				})
				m.Col(6, func() {
					m.Text(title, props.Text{
						Size:  cfg.TitleSize - 1,
						Style: consts.Bold,
						Align: consts.Right,
						Color: getPrimaryColor(),
					})
					topOffset := 9.0
					for _, line := range wideLines {
						m.Text(line, props.Text{
							Top:   topOffset,
							Size:  cfg.BodySize,
							Align: consts.Right,
							Color: getSecondaryColor(),
						})
						topOffset += 6.0
					}
				})
			})
		}
	})
}

func getBrandLogoB64() string {
	b, err := ui.Files.ReadFile("asset/img/brand_logo.png")
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func buildFooter(m pdf.Maroto, cfg paperConfig) {
	logoB64 := getBrandLogoB64()
	m.RegisterFooter(func() {
		m.Row(8.0, func() {
			if cfg.IsNarrow {
				m.Col(12, func() {
					m.Text("Powered by Teks-Invoice", props.Text{
						Top:   2,
						Size:  cfg.SmallSize,
						Style: consts.Italic,
						Align: consts.Center,
						Color: getSecondaryColor(),
					})
				})
			} else {
				m.Col(4, func() {})
				m.Col(1, func() {
					if logoB64 != "" {
						_ = m.Base64Image(logoB64, consts.Png, props.Rect{
							Percent: 80,
							Center:  true,
						})
					}
				})
				m.Col(3, func() {
					m.Text("Powered by Teks-Invoice", props.Text{
						Top:   2,
						Size:  cfg.SmallSize,
						Style: consts.Italic,
						Align: consts.Left,
						Color: getSecondaryColor(),
					})
				})
				m.Col(4, func() {})
			}
		})
	})
}
