package pdf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"raven.go.invoice-builder/internal/currency"
	"raven.go.invoice-builder/internal/models"
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

// GenerateInvoicePDF builds a styled PDF document for a given invoice, business profile, and paper format.
func GenerateInvoicePDF(invoice *models.Invoice, profile *models.BusinessProfile, paperSizeStr string) ([]byte, error) {
	m, cfg := getPaperConfig(paperSizeStr)

	// Resolve business logo file path if uploaded
	var logoPath string
	if profile != nil && profile.LogoURL != "" {
		localPath := filepath.Join(".", "ui", profile.LogoURL)
		if _, err := os.Stat(localPath); err == nil {
			logoPath = localPath
		}
	}

	companyName := "Teks-Invoice Professional"
	if profile != nil && profile.CompanyName != "" {
		companyName = profile.CompanyName
	}

	// Header Section
	m.RegisterHeader(func() {
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
			// Stacked Header layout for narrow receipts (POS 80 / POS 58)
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
					m.Text(fmt.Sprintf("INVOICE #%s", invoice.InvoiceNumber), props.Text{
						Size:  cfg.TitleSize,
						Style: consts.Bold,
						Align: consts.Center,
						Color: getPrimaryColor(),
					})
				})
			})
			m.Row(8.0, func() {
				m.Col(12, func() {
					datesStr := fmt.Sprintf("Date: %s | Due: %s",
						invoice.IssueDate.Format("Jan 02, 2006"),
						invoice.DueDate.Format("Jan 02, 2006"))
					m.Text(datesStr, props.Text{
						Size:  cfg.SmallSize,
						Align: consts.Center,
						Color: getSecondaryColor(),
					})
				})
			})
		} else {
			// Standard 2-column layout for A4, Letter, Legal, A5
			m.Row(22.0, func() {
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
					m.Text(fmt.Sprintf("INVOICE #%s", invoice.InvoiceNumber), props.Text{
						Size:  cfg.TitleSize,
						Style: consts.Bold,
						Align: consts.Right,
						Color: getPrimaryColor(),
					})
					m.Text(fmt.Sprintf("Issue Date: %s", invoice.IssueDate.Format("Jan 02, 2006")), props.Text{
						Top:   9,
						Size:  cfg.BodySize,
						Align: consts.Right,
						Color: getSecondaryColor(),
					})
					m.Text(fmt.Sprintf("Due Date: %s", invoice.DueDate.Format("Jan 02, 2006")), props.Text{
						Top:   15,
						Size:  cfg.BodySize,
						Align: consts.Right,
						Color: getSecondaryColor(),
					})
				})
			})
		}
	})

	// Main Body Content
	m.Line(6.0)

	// Billed To Section
	billedHeight := 18.0
	if cfg.IsNarrow {
		billedHeight = 12.0
	}
	m.Row(billedHeight, func() {
		m.Col(12, func() {
			clientName := "Direct Client"
			clientEmail := ""
			clientAddress := ""

			if invoice.Client != nil {
				clientName = invoice.Client.Name
				clientEmail = invoice.Client.Email
				clientAddress = invoice.Client.Address
			}

			if cfg.IsNarrow {
				m.Text(fmt.Sprintf("BILLED TO: %s", clientName), props.Text{
					Size:  cfg.BodySize,
					Style: consts.Bold,
					Color: getDarkColor(),
				})
				if clientEmail != "" {
					m.Text(clientEmail, props.Text{
						Top:   4,
						Size:  cfg.SmallSize,
						Color: getSecondaryColor(),
					})
				}
			} else {
				m.Text("BILLED TO:", props.Text{
					Size:  cfg.BodySize,
					Style: consts.Bold,
					Color: getSecondaryColor(),
				})
				m.Text(clientName, props.Text{
					Top:   5,
					Size:  cfg.TitleSize - 5,
					Style: consts.Bold,
					Color: getDarkColor(),
				})
				if clientEmail != "" {
					m.Text(clientEmail, props.Text{
						Top:   10,
						Size:  cfg.BodySize,
						Color: getSecondaryColor(),
					})
				}
				if clientAddress != "" {
					m.Text(clientAddress, props.Text{
						Top:   15,
						Size:  cfg.BodySize,
						Color: getSecondaryColor(),
					})
				}
			}
		})
	})

	m.Line(6.0)

	// Line Items Table with Responsive Row Padding
	sym := currency.Symbol(invoice.Currency)
	priceHeader := fmt.Sprintf("Unit Price (%s)", sym)
	amountHeader := fmt.Sprintf("Amount (%s)", sym)
	if cfg.IsNarrow {
		priceHeader = "Price"
		amountHeader = "Total"
	}

	m.Row(10.0, func() {
		m.Col(cfg.DescCol, func() {
			m.Text("Description", props.Text{Size: cfg.BodySize, Style: consts.Bold, Color: getDarkColor()})
		})
		m.Col(cfg.QtyCol, func() {
			m.Text("Qty", props.Text{Size: cfg.BodySize, Style: consts.Bold, Align: consts.Center, Color: getDarkColor()})
		})
		m.Col(cfg.PriceCol, func() {
			m.Text(priceHeader, props.Text{Size: cfg.BodySize, Style: consts.Bold, Align: consts.Right, Color: getDarkColor()})
		})
		m.Col(cfg.AmountCol, func() {
			m.Text(amountHeader, props.Text{Size: cfg.BodySize, Style: consts.Bold, Align: consts.Right, Color: getDarkColor()})
		})
	})
	m.Line(4.0)

	rowHeight := 12.0
	if cfg.IsNarrow {
		rowHeight = 9.0
	}

	for _, item := range invoice.LineItems {
		m.Row(rowHeight, func() {
			m.Col(cfg.DescCol, func() {
				m.Text(item.Description, props.Text{Top: 2, Size: cfg.BodySize, Color: getDarkColor()})
			})
			m.Col(cfg.QtyCol, func() {
				m.Text(fmt.Sprintf("%.2f", item.Quantity), props.Text{Top: 2, Size: cfg.BodySize, Align: consts.Center, Color: getDarkColor()})
			})
			m.Col(cfg.PriceCol, func() {
				m.Text(fmt.Sprintf("%.2f", item.UnitPrice), props.Text{Top: 2, Size: cfg.BodySize, Align: consts.Right, Color: getDarkColor()})
			})
			m.Col(cfg.AmountCol, func() {
				m.Text(fmt.Sprintf("%.2f", item.Amount), props.Text{Top: 2, Size: cfg.BodySize, Style: consts.Bold, Align: consts.Right, Color: getDarkColor()})
			})
		})
		m.Line(3.0)
	}

	m.Line(4.0)

	// Calculation Summary Section
	if cfg.IsNarrow {
		m.Row(24.0, func() {
			m.Col(12, func() {
				m.Text(fmt.Sprintf("Subtotal: %s%.2f", sym, invoice.Subtotal), props.Text{
					Size:  cfg.BodySize,
					Align: consts.Right,
					Color: getDarkColor(),
				})
				topOffset := 5.0
				if invoice.Tax > 0 {
					m.Text(fmt.Sprintf("Tax: %s%.2f", sym, invoice.Tax), props.Text{
						Top:   topOffset,
						Size:  cfg.BodySize,
						Align: consts.Right,
						Color: getDarkColor(),
					})
					topOffset += 5.0
				}
				if invoice.Discount > 0 {
					m.Text(fmt.Sprintf("Discount: -%s%.2f", sym, invoice.Discount), props.Text{
						Top:   topOffset,
						Size:  cfg.BodySize,
						Align: consts.Right,
						Color: getDarkColor(),
					})
					topOffset += 5.0
				}
				m.Text(fmt.Sprintf("TOTAL DUE: %s%.2f", sym, invoice.Total), props.Text{
					Top:   topOffset,
					Size:  cfg.TitleSize,
					Style: consts.Bold,
					Align: consts.Right,
					Color: getPrimaryColor(),
				})
			})
		})
		if invoice.Notes != "" {
			m.Row(14.0, func() {
				m.Col(12, func() {
					m.Text("Notes & Terms:", props.Text{
						Size:  cfg.SmallSize,
						Style: consts.Bold,
						Color: getSecondaryColor(),
					})
					m.Text(invoice.Notes, props.Text{
						Top:   4,
						Size:  cfg.SmallSize,
						Color: getDarkColor(),
					})
				})
			})
		}
	} else {
		m.Row(28.0, func() {
			m.Col(6, func() {
				if invoice.Notes != "" {
					m.Text("Notes & Terms:", props.Text{
						Size:  cfg.BodySize,
						Style: consts.Bold,
						Color: getSecondaryColor(),
					})
					m.Text(invoice.Notes, props.Text{
						Top:   5,
						Size:  cfg.SmallSize,
						Color: getDarkColor(),
					})
				}
			})
			m.Col(6, func() {
				m.Text(fmt.Sprintf("Subtotal: %s%.2f", sym, invoice.Subtotal), props.Text{
					Size:  cfg.BodySize,
					Align: consts.Right,
					Color: getDarkColor(),
				})
				if invoice.Tax > 0 {
					m.Text(fmt.Sprintf("Tax: %s%.2f", sym, invoice.Tax), props.Text{
						Top:   5,
						Size:  cfg.BodySize,
						Align: consts.Right,
						Color: getDarkColor(),
					})
				}
				if invoice.Discount > 0 {
					m.Text(fmt.Sprintf("Discount: -%s%.2f", sym, invoice.Discount), props.Text{
						Top:   10,
						Size:  cfg.BodySize,
						Align: consts.Right,
						Color: getDarkColor(),
					})
				}
				m.Text(fmt.Sprintf("TOTAL DUE: %s%.2f", sym, invoice.Total), props.Text{
					Top:   16,
					Size:  cfg.TitleSize - 3,
					Style: consts.Bold,
					Align: consts.Right,
					Color: getPrimaryColor(),
				})
			})
		})
	}

	// Branded footer
	appLogoPath := filepath.Join(".", "ui", "asset", "img", "brand_logo.png")
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
					if _, err := os.Stat(appLogoPath); err == nil {
						_ = m.FileImage(appLogoPath, props.Rect{
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

	buffer, err := m.Output()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
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

// GenerateReceiptPDF builds a styled PDF payment receipt document for a given paper size.
func GenerateReceiptPDF(receipt *models.Receipt, invoice *models.Invoice, profile *models.BusinessProfile, paperSizeStr string) ([]byte, error) {
	m, cfg := getPaperConfig(paperSizeStr)

	var logoPath string
	if profile != nil && profile.LogoURL != "" {
		localPath := filepath.Join(".", "ui", profile.LogoURL)
		if _, err := os.Stat(localPath); err == nil {
			logoPath = localPath
		}
	}

	companyName := "Teks-Invoice Professional"
	if profile != nil && profile.CompanyName != "" {
		companyName = profile.CompanyName
	}

	m.RegisterHeader(func() {
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
			m.Row(8.0, func() {
				m.Col(12, func() {
					m.Text(fmt.Sprintf("RECEIPT #%s", receipt.ReceiptNumber), props.Text{
						Size:  cfg.TitleSize - 2,
						Style: consts.Bold,
						Align: consts.Center,
						Color: getPrimaryColor(),
					})
				})
			})
			m.Row(8.0, func() {
				m.Col(12, func() {
					m.Text(fmt.Sprintf("Date: %s | Inv #: %s", receipt.IssuedAt.Format("Jan 02, 2006"), invoice.InvoiceNumber), props.Text{
						Size:  cfg.SmallSize,
						Align: consts.Center,
						Color: getSecondaryColor(),
					})
				})
			})
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
							Top:   7,
							Size:  cfg.BodySize,
							Color: getSecondaryColor(),
						})
					}
					if profile != nil && profile.RegistrationNumber != "" {
						m.Text(fmt.Sprintf("Reg #: %s", profile.RegistrationNumber), props.Text{
							Top:   13,
							Size:  cfg.SmallSize,
							Color: getSecondaryColor(),
						})
					}
				})
				m.Col(6, func() {
					m.Text("PAYMENT RECEIPT", props.Text{
						Size:  cfg.TitleSize - 3,
						Style: consts.Bold,
						Align: consts.Right,
						Color: getPrimaryColor(),
					})
					m.Text(fmt.Sprintf("Receipt #: %s", receipt.ReceiptNumber), props.Text{
						Top:   7,
						Size:  cfg.BodySize,
						Style: consts.Bold,
						Align: consts.Right,
						Color: getDarkColor(),
					})
					m.Text(fmt.Sprintf("Date Issued: %s", receipt.IssuedAt.Format("Jan 02, 2006")), props.Text{
						Top:   13,
						Size:  cfg.BodySize,
						Align: consts.Right,
						Color: getSecondaryColor(),
					})
					m.Text(fmt.Sprintf("For Invoice #: %s", invoice.InvoiceNumber), props.Text{
						Top:   19,
						Size:  cfg.BodySize,
						Align: consts.Right,
						Color: getSecondaryColor(),
					})
				})
			})
		}
	})

	m.Line(6.0)

	// Received From / Client Section
	billedHeight := 18.0
	if cfg.IsNarrow {
		billedHeight = 12.0
	}
	m.Row(billedHeight, func() {
		m.Col(12, func() {
			clientName := "Direct Client"
			clientEmail := ""
			if invoice.Client != nil {
				clientName = invoice.Client.Name
				clientEmail = invoice.Client.Email
			}

			if cfg.IsNarrow {
				m.Text(fmt.Sprintf("RECEIVED FROM: %s", clientName), props.Text{
					Size:  cfg.BodySize,
					Style: consts.Bold,
					Color: getDarkColor(),
				})
			} else {
				m.Text("RECEIVED FROM:", props.Text{
					Size:  cfg.BodySize,
					Style: consts.Bold,
					Color: getSecondaryColor(),
				})
				m.Text(clientName, props.Text{
					Top:   5,
					Size:  cfg.TitleSize - 5,
					Style: consts.Bold,
					Color: getDarkColor(),
				})
				if clientEmail != "" {
					m.Text(clientEmail, props.Text{
						Top:   10,
						Size:  cfg.BodySize,
						Color: getSecondaryColor(),
					})
				}
			}
		})
	})

	m.Line(6.0)

	// Itemized Breakdown Table
	sym := currency.Symbol(receipt.Currency)
	priceHeader := fmt.Sprintf("Unit Price (%s)", sym)
	amountHeader := fmt.Sprintf("Amount (%s)", sym)
	if cfg.IsNarrow {
		priceHeader = "Price"
		amountHeader = "Total"
	}

	m.Row(10.0, func() {
		m.Col(cfg.DescCol, func() {
			m.Text("Item Description", props.Text{Size: cfg.BodySize, Style: consts.Bold, Color: getDarkColor()})
		})
		m.Col(cfg.QtyCol, func() {
			m.Text("Qty", props.Text{Size: cfg.BodySize, Style: consts.Bold, Align: consts.Center, Color: getDarkColor()})
		})
		m.Col(cfg.PriceCol, func() {
			m.Text(priceHeader, props.Text{Size: cfg.BodySize, Style: consts.Bold, Align: consts.Right, Color: getDarkColor()})
		})
		m.Col(cfg.AmountCol, func() {
			m.Text(amountHeader, props.Text{Size: cfg.BodySize, Style: consts.Bold, Align: consts.Right, Color: getDarkColor()})
		})
	})
	m.Line(4.0)

	rowHeight := 12.0
	if cfg.IsNarrow {
		rowHeight = 9.0
	}

	if len(invoice.LineItems) > 0 {
		for _, item := range invoice.LineItems {
			m.Row(rowHeight, func() {
				m.Col(cfg.DescCol, func() {
					m.Text(item.Description, props.Text{Top: 2, Size: cfg.BodySize, Color: getDarkColor()})
				})
				m.Col(cfg.QtyCol, func() {
					m.Text(fmt.Sprintf("%.2f", item.Quantity), props.Text{Top: 2, Size: cfg.BodySize, Align: consts.Center, Color: getDarkColor()})
				})
				m.Col(cfg.PriceCol, func() {
					m.Text(fmt.Sprintf("%.2f", item.UnitPrice), props.Text{Top: 2, Size: cfg.BodySize, Align: consts.Right, Color: getDarkColor()})
				})
				m.Col(cfg.AmountCol, func() {
					m.Text(fmt.Sprintf("%.2f", item.Amount), props.Text{Top: 2, Size: cfg.BodySize, Style: consts.Bold, Align: consts.Right, Color: getDarkColor()})
				})
			})
			m.Line(3.0)
		}
	} else {
		m.Row(rowHeight, func() {
			m.Col(6, func() {
				m.Text(fmt.Sprintf("Payment Settlement for Invoice %s", invoice.InvoiceNumber), props.Text{Top: 3, Size: cfg.BodySize, Color: getDarkColor()})
			})
			m.Col(6, func() {
				m.Text(fmt.Sprintf("%s%.2f", sym, receipt.Amount), props.Text{Top: 3, Size: cfg.TitleSize - 3, Style: consts.Bold, Align: consts.Right, Color: getPrimaryColor()})
			})
		})
		m.Line(3.0)
	}

	m.Line(4.0)

	// Summary Section
	m.Row(24.0, func() {
		m.Col(6, func() {
			m.Text("Status: PAID IN FULL", props.Text{
				Size:  cfg.BodySize,
				Style: consts.Bold,
				Color: getPrimaryColor(),
			})
		})
		m.Col(6, func() {
			m.Text(fmt.Sprintf("Subtotal: %s%.2f", sym, invoice.Subtotal), props.Text{
				Size:  cfg.BodySize,
				Align: consts.Right,
				Color: getDarkColor(),
			})
			topOffset := 5.0
			if invoice.Tax > 0 {
				m.Text(fmt.Sprintf("Tax: %s%.2f", sym, invoice.Tax), props.Text{
					Top:   topOffset,
					Size:  cfg.BodySize,
					Align: consts.Right,
					Color: getDarkColor(),
				})
				topOffset += 5.0
			}
			if invoice.Discount > 0 {
				m.Text(fmt.Sprintf("Discount: -%s%.2f", sym, invoice.Discount), props.Text{
					Top:   topOffset,
					Size:  cfg.BodySize,
					Align: consts.Right,
					Color: getDarkColor(),
				})
				topOffset += 5.0
			}
			m.Text(fmt.Sprintf("TOTAL PAID: %s%.2f", sym, receipt.Amount), props.Text{
				Top:   topOffset,
				Size:  cfg.TitleSize - 3,
				Style: consts.Bold,
				Align: consts.Right,
				Color: getPrimaryColor(),
			})
		})
	})

	m.RegisterFooter(func() {
		m.Row(8.0, func() {
			m.Col(12, func() {
				m.Text("Official Payment Receipt - Generated by Teks-Invoice", props.Text{
					Top:   2,
					Size:  cfg.SmallSize,
					Align: consts.Center,
					Color: getSecondaryColor(),
				})
			})
		})
	})

	buffer, err := m.Output()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}