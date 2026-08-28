package pdf

import (
	"testing"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
)

// TestMarotoEngineCast ensures that the *pdf.PdfMaroto cast we rely on to bypass
// the grid layout engine for the watermark injection remains valid.
func TestMarotoEngineCast(t *testing.T) {
	// Initialize a standard Maroto instance
	m := pdf.NewMaroto(consts.Portrait, consts.A4)

	// Attempt to cast to PdfMaroto
	_, ok := m.(*pdf.PdfMaroto)
	if !ok {
		t.Fatalf("CRITICAL: Failed to cast pdf.Maroto to *pdf.PdfMaroto. This means the underlying gofpdf engine is no longer accessible via this struct, which breaks the watermark injection logic. If Maroto was recently upgraded, you must find a new way to access gofpdf or inject absolute-positioned images.")
	}
}
