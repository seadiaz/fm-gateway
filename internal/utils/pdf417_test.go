package utils

import (
	"strings"
	"testing"
)

func TestGenerateStampPDF417(t *testing.T) {
	// Create test stamp data
	testCAF := PDF417CAF{
		Version: "1.0",
		DA: PDF417DA{
			RE:  "12345678-9",
			RS:  "Test Company S.A.",
			TD:  33,
			RNG: PDF417RNG{D: 1, H: 100},
			FA:  "2024-01-15",
			RSAPK: PDF417RSAPK{
				M: "test_modulus",
				E: "test_exponent",
			},
			IDK: "test_idk",
		},
		FRMA: PDF417FRMA{
			Algorithm: "SHA1withRSA",
			Value:     "test_signature",
		},
	}

	testStampData := StampPDF417Data{
		Version: "1.0",
		DD: PDF417DD{
			RE:    "12345678-9",
			TD:    33,
			F:     1,
			FE:    "2024-01-15",
			RR:    "11111111-1",
			RSR:   "Test Client",
			MNT:   20000,
			IT1:   "Test Product",
			CAF:   testCAF,
			TSTED: "2024-01-15T10:30:00",
		},
		FRMT: PDF417FRMT{
			Algorithm: "SHA1withRSA",
			Value:     "test_stamp_signature",
		},
	}

	// Generate PDF417 barcode
	img, pngBytes, err := GenerateStampPDF417(testStampData, 350, 80)
	if err != nil {
		t.Fatalf("GenerateStampPDF417 failed: %v", err)
	}

	// Verify image was created
	if img == nil {
		t.Error("Generated image should not be nil")
	}

	// Verify PNG bytes were created
	if len(pngBytes) == 0 {
		t.Error("PNG bytes should not be empty")
	}

	// Verify image dimensions
	bounds := img.Bounds()
	if bounds.Dx() != 350 || bounds.Dy() != 80 {
		t.Errorf("Expected dimensions 350x80, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	t.Logf("Successfully generated PDF417 barcode with %d bytes", len(pngBytes))
}

func TestGenerateStampPDF417FromXML(t *testing.T) {
	// Test XML data (simplified TED structure)
	testXML := `<TED version="1.0">
	<DD>
		<RE>12345678-9</RE>
		<TD>33</TD>
		<F>1</F>
		<FE>2024-01-15</FE>
		<RR>11111111-1</RR>
		<RSR>Test Client</RSR>
		<MNT>20000</MNT>
		<IT1>Test Product</IT1>
		<TSTED>2024-01-15T10:30:00</TSTED>
	</DD>
	<FRMT algoritmo="SHA1withRSA">test_signature</FRMT>
</TED>`

	// Generate PDF417 barcode from XML
	img, pngBytes, err := GenerateStampPDF417FromXML(testXML, 300, 120)
	if err != nil {
		t.Fatalf("GenerateStampPDF417FromXML failed: %v", err)
	}

	// Verify image was created
	if img == nil {
		t.Error("Generated image should not be nil")
	}

	// Verify PNG bytes were created
	if len(pngBytes) == 0 {
		t.Error("PNG bytes should not be empty")
	}

	// Verify image dimensions
	bounds := img.Bounds()
	if bounds.Dx() != 300 || bounds.Dy() != 120 {
		t.Errorf("Expected dimensions 300x120, got %dx%d", bounds.Dx(), bounds.Dy())
	}

	t.Logf("Successfully generated PDF417 barcode from XML with %d bytes", len(pngBytes))
}

func TestConvertDomainStampToPDF417Data(t *testing.T) {
	// Test data
	testCAF := PDF417CAF{
		Version: "1.0",
		DA: PDF417DA{
			RE:  "12345678-9",
			RS:  "Test Company S.A.",
			TD:  33,
			RNG: PDF417RNG{D: 1, H: 100},
			FA:  "2024-01-15",
			RSAPK: PDF417RSAPK{
				M: "test_modulus",
				E: "test_exponent",
			},
			IDK: "test_idk",
		},
		FRMA: PDF417FRMA{
			Algorithm: "SHA1withRSA",
			Value:     "test_signature",
		},
	}

	// Convert domain data to PDF417 structure
	pdf417Data := ConvertDomainStampToPDF417Data(
		"12345678-9",           // re
		33,                     // td
		1,                      // f
		"2024-01-15",           // fe
		"11111111-1",           // rr
		"Test Client",          // rsr
		20000,                  // mnt
		"Test Product",         // it1
		testCAF,                // caf
		"2024-01-15T10:30:00",  // tsted
		"test_stamp_signature", // frmt
	)

	// Verify the conversion
	if pdf417Data.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", pdf417Data.Version)
	}

	if pdf417Data.DD.RE != "12345678-9" {
		t.Errorf("Expected RE 12345678-9, got %s", pdf417Data.DD.RE)
	}

	if pdf417Data.DD.F != 1 {
		t.Errorf("Expected F 1, got %d", pdf417Data.DD.F)
	}

	if pdf417Data.FRMT.Algorithm != "SHA1withRSA" {
		t.Errorf("Expected algorithm SHA1withRSA, got %s", pdf417Data.FRMT.Algorithm)
	}

	t.Log("Successfully converted domain stamp to PDF417 data structure")
}

func TestPDF417WithLargeData(t *testing.T) {
	// Test with larger data to ensure PDF417 can handle real-world stamp data
	largeDescription := strings.Repeat("Large Product Description ", 10)

	testXML := `<TED version="1.0">
	<DD>
		<RE>12345678-9</RE>
		<TD>33</TD>
		<F>12345</F>
		<FE>2024-01-15</FE>
		<RR>11111111-1</RR>
		<RSR>` + largeDescription + `</RSR>
		<MNT>999999999</MNT>
		<IT1>` + largeDescription + `</IT1>
		<TSTED>2024-01-15T10:30:00-03:00</TSTED>
	</DD>
	<FRMT algoritmo="SHA1withRSA">very_long_signature_` + strings.Repeat("x", 100) + `</FRMT>
</TED>`

	// Generate PDF417 barcode with larger dimensions for complex data
	// Using dimensions that meet the minimum requirements (at least 409x56)
	img, pngBytes, err := GenerateStampPDF417FromXML(testXML, 500, 200)
	if err != nil {
		t.Fatalf("GenerateStampPDF417FromXML with large data failed: %v", err)
	}

	// Verify the barcode was generated successfully
	if img == nil {
		t.Error("Generated image should not be nil")
	}

	if len(pngBytes) == 0 {
		t.Error("PNG bytes should not be empty")
	}

	t.Logf("Successfully generated PDF417 barcode for large data with %d bytes", len(pngBytes))
}

func TestCalculateSafePDF417Dimensions(t *testing.T) {
	tests := []struct {
		name       string
		dataLength int
		minWidth   int
		minHeight  int
	}{
		{"Small data", 500, 400, 150},
		{"Medium data", 1200, 450, 180},
		{"Standard data", 1800, 500, 200},
		{"Large data", 2500, 550, 220},
		{"Very large data", 3000, 600, 250},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			width, height := CalculateSafePDF417Dimensions(tt.dataLength)

			if width < tt.minWidth {
				t.Errorf("Width %d is less than minimum expected %d for %s", width, tt.minWidth, tt.name)
			}

			if height < tt.minHeight {
				t.Errorf("Height %d is less than minimum expected %d for %s", height, tt.minHeight, tt.name)
			}

			t.Logf("%s: %d bytes → %dx%d dimensions", tt.name, tt.dataLength, width, height)
		})
	}
}

func TestGenerateStampPDF417FromXMLWithAutoDimensions(t *testing.T) {
	// Test with realistic TED XML (PDF417 has ~1800 character limit)
	realWorldXML := `<TED version="1.0">
	<DD>
		<RE>76543210-9</RE>
		<TD>33</TD>
		<F>98765</F>
		<FE>2024-01-15</FE>
		<RR>11222333-4</RR>
		<RSR>Cliente Empresa Larga S.A.</RSR>
		<MNT>1234567890</MNT>
		<IT1>Producto con Descripción Detallada</IT1>
		<CAF version="1.0">
			<DA>
				<RE>76543210-9</RE>
				<RS>Mi Empresa Consultora S.A.</RS>
				<TD>33</TD>
				<RNG>
					<D>1</D>
					<H>50000</H>
				</RNG>
				<FA>2024-01-15</FA>
				<RSAPK>
					<M>nZABIUCtdK4Fn1J8qY8MiE2DRPaCsG/qGgzFpJAuCBU=</M>
					<E>AQAB</E>
				</RSAPK>
				<IDK>12345</IDK>
			</DA>
			<FRMA algoritmo="SHA1withRSA">cyCgbgAAJ2D1ZR5FqON1ZmXU2CgPQTERl0LMcPnOJMeDMOlLYE7UtOwF2g==</FRMA>
		</CAF>
		<TSTED>2024-01-15T10:30:45-03:00</TSTED>
	</DD>
	<FRMT algoritmo="SHA1withRSA">JKiGxvxmprOl+8ZixIYRJVL/2Og+CY3at7W/IYEW17q/UjCd3c4r09czi1o07kLMFSdKQ8RCM9jpsWsFHPKjy1V96Eo2rf73QQ/dkPCZ5zXxFGGdN7pxqM2Bv+2Y+FzrfGbSpmkbMTVdEbNlEFgeANNVyPeVOoPy5gXo3YtUgFw=</FRMT>
</TED>`

	// Generate barcode with auto dimensions
	img, pngBytes, err := GenerateStampPDF417FromXMLWithAutoDimensions(realWorldXML)
	if err != nil {
		t.Fatalf("GenerateStampPDF417FromXMLWithAutoDimensions failed: %v", err)
	}

	// Verify the barcode was generated successfully
	if img == nil {
		t.Error("Generated image should not be nil")
	}

	if len(pngBytes) == 0 {
		t.Error("PNG bytes should not be empty")
	}

	bounds := img.Bounds()
	t.Logf("Auto-dimensions barcode: %dx%d pixels, %d bytes PNG, %d bytes XML",
		bounds.Dx(), bounds.Dy(), len(pngBytes), len(realWorldXML))

	// Verify dimensions are reasonable for the data size
	if bounds.Dx() < 500 || bounds.Dy() < 200 {
		t.Errorf("Dimensions %dx%d seem too small for data size %d", bounds.Dx(), bounds.Dy(), len(realWorldXML))
	}
}

func TestGenerateStampPDF417WithAutoDimensions(t *testing.T) {
	// Create a realistic stamp with appropriate data sizes for PDF417
	complexCAF := PDF417CAF{
		Version: "1.0",
		DA: PDF417DA{
			RE:  "76543210-9",
			RS:  "Mi Empresa Consultora S.A.",
			TD:  33,
			RNG: PDF417RNG{D: 1, H: 50000},
			FA:  "2024-01-15",
			RSAPK: PDF417RSAPK{
				M: "nZABIUCtdK4Fn1J8qY8MiE2DRPaCsG/qGgzFpJAuCBU=",
				E: "AQAB",
			},
			IDK: "12345",
		},
		FRMA: PDF417FRMA{
			Algorithm: "SHA1withRSA",
			Value:     "cyCgbgAAJ2D1ZR5FqON1ZmXU2CgPQTERl0LMcPnOJMeDMOlLYE7UtOwF2g==",
		},
	}

	complexStampData := StampPDF417Data{
		Version: "1.0",
		DD: PDF417DD{
			RE:    "76543210-9",
			TD:    33,
			F:     98765,
			FE:    "2024-01-15",
			RR:    "11222333-4",
			RSR:   "Cliente Empresa Larga S.A.",
			MNT:   1234567890,
			IT1:   "Producto con Descripción Detallada",
			CAF:   complexCAF,
			TSTED: "2024-01-15T10:30:45-03:00",
		},
		FRMT: PDF417FRMT{
			Algorithm: "SHA1withRSA",
			Value:     "JKiGxvxmprOl+8ZixIYRJVL/2Og+CY3at7W/IYEW17q/UjCd3c4r09czi1o07kLMFSdKQ8RCM9jpsWsFHPKjy1V96Eo2rf73QQ/dkPCZ5zXxFGGdN7pxqM2Bv+2Y+FzrfGbSpmkbMTVdEbNlEFgeANNVyPeVOoPy5gXo3YtUgFw=",
		},
	}

	// Generate barcode with auto dimensions
	img, pngBytes, err := GenerateStampPDF417WithAutoDimensions(complexStampData)
	if err != nil {
		t.Fatalf("GenerateStampPDF417WithAutoDimensions failed: %v", err)
	}

	// Verify the barcode was generated successfully
	if img == nil {
		t.Error("Generated image should not be nil")
	}

	if len(pngBytes) == 0 {
		t.Error("PNG bytes should not be empty")
	}

	bounds := img.Bounds()
	t.Logf("Complex auto-dimensions barcode: %dx%d pixels, %d bytes PNG",
		bounds.Dx(), bounds.Dy(), len(pngBytes))
}
