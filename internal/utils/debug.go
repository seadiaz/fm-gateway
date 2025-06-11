package utils

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/net/html/charset"
)

// TestPrivateKeyFromCAF reads a CAF file and tests private key parsing
func TestPrivateKeyFromCAF(cafFilePath string) error {
	data, err := os.ReadFile(cafFilePath)
	if err != nil {
		return fmt.Errorf("reading CAF file: %w", err)
	}

	var cafXML struct {
		XMLName xml.Name `xml:"AUTORIZACION"`
		CAF     struct {
			DA struct {
				RE string `xml:"RE"`
			} `xml:"DA"`
		} `xml:"CAF"`
		RSASK struct {
			Value string `xml:",chardata"`
		} `xml:"RSASK"`
	}

	// Use the same decoder as the controller
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&cafXML)
	if err != nil {
		return fmt.Errorf("parsing XML: %w", err)
	}

	privateKey := cafXML.RSASK.Value
	slog.Info("Debug: Private key from CAF",
		slog.String("company", cafXML.CAF.DA.RE),
		slog.Int("keyLength", len(privateKey)),
		slog.String("keyStart", privateKey[:50]),
		slog.String("keyEnd", privateKey[len(privateKey)-50:]))

	// Test signing
	testData := []byte("test data")
	signature, err := SignSHA1WithRSA(testData, privateKey)
	if err != nil {
		return fmt.Errorf("testing signature: %w", err)
	}

	slog.Info("Success! Generated signature", slog.String("signature", signature))
	return nil
}
