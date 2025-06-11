package utils

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"fmt"
	"strings"
)

// SignSHA1WithRSA creates a SHA1withRSA digital signature
func SignSHA1WithRSA(data []byte, privateKeyPEM string) (string, error) {
	if len(privateKeyPEM) == 0 {
		return "", fmt.Errorf("private key is empty")
	}

	// Clean and normalize the private key to proper PEM format
	cleanedKey := cleanPrivateKey(privateKeyPEM)

	// Parse the private key from PEM format
	block, _ := pem.Decode([]byte(cleanedKey))
	if block == nil {
		return "", fmt.Errorf("failed to decode PEM block containing private key (key length: %d, starts with: %.50s)",
			len(cleanedKey), cleanedKey)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8 format as well
		pkcs8Key, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err2 != nil {
			return "", fmt.Errorf("failed to parse private key: PKCS1 error: %w, PKCS8 error: %v", err, err2)
		}

		var ok bool
		privateKey, ok = pkcs8Key.(*rsa.PrivateKey)
		if !ok {
			return "", fmt.Errorf("PKCS8 key is not an RSA key")
		}
	}

	// Create SHA1 hash of the data
	hash := sha1.Sum(data)

	// Sign the hash using RSA-PSS
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %w", err)
	}

	// Encode the signature to base64
	return base64.StdEncoding.EncodeToString(signature), nil
}

// SerializeToXMLWithoutNewlines serializes a struct to XML and removes all newlines
func SerializeToXMLWithoutNewlines(v interface{}) ([]byte, error) {
	xmlData, err := xml.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal to XML: %w", err)
	}

	// Remove all newlines and extra whitespace
	xmlString := strings.ReplaceAll(string(xmlData), "\n", "")
	xmlString = strings.ReplaceAll(xmlString, "\r", "")
	xmlString = strings.ReplaceAll(xmlString, "\t", "")

	// Remove extra spaces between tags
	xmlString = strings.TrimSpace(xmlString)

	return []byte(xmlString), nil
}

// cleanPrivateKey normalizes a private key string to proper PEM format
func cleanPrivateKey(privateKeyPEM string) string {
	// Remove any XML whitespace and normalize line endings
	cleaned := strings.TrimSpace(privateKeyPEM)

	// Split into lines and process
	lines := strings.Split(cleaned, "\n")
	var normalizedLines []string
	var isInData bool

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if strings.Contains(trimmed, "BEGIN RSA PRIVATE KEY") {
			normalizedLines = append(normalizedLines, "-----BEGIN RSA PRIVATE KEY-----")
			isInData = true
			continue
		}

		if strings.Contains(trimmed, "END RSA PRIVATE KEY") {
			normalizedLines = append(normalizedLines, "-----END RSA PRIVATE KEY-----")
			isInData = false
			continue
		}

		if isInData {
			// For data lines, we need to ensure they're properly formatted
			// Remove any embedded newlines or spaces within the base64 data
			cleanData := strings.ReplaceAll(trimmed, " ", "")
			cleanData = strings.ReplaceAll(cleanData, "\t", "")
			cleanData = strings.ReplaceAll(cleanData, "\r", "")
			if cleanData != "" {
				normalizedLines = append(normalizedLines, cleanData)
			}
		} else {
			normalizedLines = append(normalizedLines, trimmed)
		}
	}

	return strings.Join(normalizedLines, "\n")
}
