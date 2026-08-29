package pdf

import (
	"fmt"
	_ "image/jpeg"
	"os"
	"path/filepath"

	"github.com/johnfercher/maroto/pkg/consts"

	"github.com/johnfercher/maroto/pkg/props"

	"raven.go.invoice-builder/internal/currency"
	"raven.go.invoice-builder/internal/models"
)


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
	title := fmt.Sprintf("INVOICE #%s", invoice.InvoiceNumber)
	var wideLines []string
	var narrowLines []string
	if cfg.IsNarrow {
		narrowLines = []string{
			fmt.Sprintf("Date: %s | Due: %s", invoice.IssueDate.Format("Jan 02, 2006"), invoice.DueDate.Format("Jan 02, 2006")),
		}
	} else {
		wideLines = []string{
			fmt.Sprintf("Issue Date: %s", invoice.IssueDate.Format("Jan 02, 2006")),
			fmt.Sprintf("Due Date: %s", invoice.DueDate.Format("Jan 02, 2006")),
		}
	}
	buildHeader(m, cfg, profile, logoPath, companyName, title, wideLines, narrowLines)

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
	buildFooter(m, cfg)

	buffer, err := m.Output()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
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

	var title string
	var wideLines []string
	var narrowLines []string
	if cfg.IsNarrow {
		title = fmt.Sprintf("RECEIPT #%s", receipt.ReceiptNumber)
		narrowLines = []string{
			fmt.Sprintf("Date: %s | Inv #: %s", receipt.IssuedAt.Format("Jan 02, 2006"), invoice.InvoiceNumber),
		}
	} else {
		title = "PAYMENT RECEIPT"
		wideLines = []string{
			fmt.Sprintf("Receipt #: %s", receipt.ReceiptNumber),
			fmt.Sprintf("Date: %s", receipt.IssuedAt.Format("Jan 02, 2006")),
			fmt.Sprintf("Invoice #: %s", invoice.InvoiceNumber),
		}
	}
	buildHeader(m, cfg, profile, logoPath, companyName, title, wideLines, narrowLines)

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

	buildFooter(m, cfg)

	buffer, err := m.Output()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}