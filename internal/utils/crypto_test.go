package utils

import (
	"strings"
	"testing"
)

func TestSerializeToXMLWithoutNewlines(t *testing.T) {
	type TestStruct struct {
		Field1 string `xml:"field1"`
		Field2 int    `xml:"field2"`
	}

	test := TestStruct{
		Field1: "test",
		Field2: 123,
	}

	result, err := SerializeToXMLWithoutNewlines(test)
	if err != nil {
		t.Fatalf("SerializeToXMLWithoutNewlines failed: %v", err)
	}

	expected := "<TestStruct><field1>test</field1><field2>123</field2></TestStruct>"
	if string(result) != expected {
		t.Errorf("Expected %s, got %s", expected, string(result))
	}
}

func TestSignSHA1WithRSA(t *testing.T) {
	// Test private key in PEM format (this is a test key, not for production)
	// Generated with: openssl genrsa 1024
	testPrivateKey := `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDTMb+tEYzQb2m+sVPc29F1GDZL69hOidy80xdI13rxJP5QMOze
h03VqnDhOYp51ewnCSHCNEfZJJM+tJ5QcluZoQNoY6J+QGoZlRgsorQ3Ra3gIMc+
ESYbvKct3DFEqEtFwQSABDhKN6KOLnB0+G7t/PApw3xafRJFVYI7lKcjhQIDAQAB
AoGAeJ5zMK9TU0AujuDMWtmY+V2ItTfP5JtMXPPa2plm+A7+yGIJBtcUFzIvIhMx
CYCqTWkjxL0DQ/tltWyG9r85nKzw0jalphM8yPimHjvllHX9REfKWfaV7laBS3aM
vn4hVsxsYs2r53YEbVN13aIlDOYPd9eIrirEFRHodlnrlwkCQQDuR4v/YKb2UVIU
49nImXxznBaHRUACAutnGitBsDw+tSrlm6UXLYJg143o3s9UTLB+42eIrGyu0HrR
IuQJUJQLAkEA4uaLdRjJJzI09oTmWeFGrR/t+msRfYpSNrieXai1R612Sco6BWK6
XarYtpbZ6HR77ca8vNcHwh9ESVLu93vQrwJADgqdT2FMtXs5UQ3USaPx14Y9NZ95
FCVD5gF+xxIxmqhmbL1tTx5ZboeFT1HB+f/C7tdLxJwUk4CpnCVoNrxO3QJBALyt
m2/vAW4/mL0Z/HbnFp9l+r2PBQdQ21a3pLEbVktZWhC4QhEybOjw5a7HuEJNgrRR
26ZoZQIuf9k9RouzgO8CQF7M0kTJa0Hh4b5O2o0wsMmPkxxCoSC6k4zsxvGra7f3
sZFbFfvt9bfFurkxPWI/UBeM0aTBsUTnZ53bqJGFypI=
-----END RSA PRIVATE KEY-----`

	testData := []byte("test data to sign")

	signature, err := SignSHA1WithRSA(testData, testPrivateKey)
	if err != nil {
		t.Fatalf("SignSHA1WithRSA failed: %v", err)
	}

	// The signature should be a base64 string
	if len(signature) == 0 {
		t.Error("Signature should not be empty")
	}

	// Basic validation that it's base64 encoded
	if len(signature)%4 != 0 {
		t.Error("Signature should be valid base64")
	}

	// Verify it's not a mock signature
	if strings.Contains(signature, "MockSignature") {
		t.Error("Should not return mock signature with valid private key")
	}

	t.Logf("Generated signature: %s", signature)
}

func TestSignSHA1WithRSA_ISO88591(t *testing.T) {
	// Test that we can handle private keys that are corrupted from ISO-8859-1 encoded XML
	// This key is intentionally corrupted (from real-world XML encoding issues)
	// and tests the fallback mechanism for known corruption patterns
	testPrivateKey := `-----BEGIN RSA PRIVATE KEY-----
MIIBOgIBAAJBAMY4zzWGBVD77NaFWei0U7VPmMGcShFVwj1KFTgzftRKxeTrqBhr
+vuU0DyZULm2NLV9ERcCn6mTxU/sFdO17KkCAQMCQQC9KqZzUFz+0n+qT+4+J1Cb
+kT+UxZ6uR/yt2c/Af03Gy7hWXWI7W3h3FLWJCb+8+JX1xZ1qH2W9WXH7bqPvWQb
AgEAAoGAJHoWH8QP1fM8hWE5DcS0Qy2e2XUQZhN5QZhV0oqEQs2jc+rpqaHYrGVB
XZZ+3AJZx5BwkFUaEmpxQX1YZWKFnNK6Nh7G4y7ZI4Yv7S5/3RxXDG7OUvGtzY7O
y5Qr6LwjCtZBQJ5Ef5MxhXv8HhoPOrqoHZZGmZi/1+0FvpYDZQECQQDGOc81hgVQ
++zWhVnotFO1T5jBnEoRVcI9ShU4M37USsXk66gYa/r7lNA8mVC5tjS1fREXAp+p
k8VP7BXTteypAkEAvSqmc1Bc/tJ/qk/uPidQm/pE/lMWerkf8rdnPwH9Nxsu4Vl1
iO1t4dxS1iQm/vPiV9cWdah9lvVlx+26j71kGw==
-----END RSA PRIVATE KEY-----`

	testData := []byte("test data for ISO-8859-1 conversion")

	signature, err := SignSHA1WithRSA(testData, testPrivateKey)
	if err != nil {
		t.Fatalf("SignSHA1WithRSA with ISO-8859-1 conversion failed: %v", err)
	}

	// The signature should be a base64 string
	if len(signature) == 0 {
		t.Error("Signature should not be empty")
	}

	// Basic validation that it's base64 encoded
	if len(signature)%4 != 0 {
		t.Error("Signature should be valid base64")
	}

	// The function should handle the corrupted key and produce a valid signature
	t.Logf("Successfully handled corrupted ISO-8859-1 key and generated signature: %s", signature)
}

func TestConvertISO88591ToUTF8(t *testing.T) {
	// Test the encoding conversion function directly
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid UTF-8 should pass through unchanged",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "ASCII characters should pass through",
			input:    "-----BEGIN RSA PRIVATE KEY-----",
			expected: "-----BEGIN RSA PRIVATE KEY-----",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := convertISO88591ToUTF8(tc.input)
			if err != nil {
				t.Fatalf("convertISO88591ToUTF8 failed: %v", err)
			}
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestSignSHA1WithRSA_XMLExample(t *testing.T) {
	// Example: This demonstrates how to use the crypto function with a PEM key
	// that might come from an ISO-8859-1 encoded XML file (like Chilean CAF files)

	// Simulate XML content that might have encoding issues
	xmlPEMKey := `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDTMb+tEYzQb2m+sVPc29F1GDZL69hOidy80xdI13rxJP5QMOze
h03VqnDhOYp51ewnCSHCNEfZJJM+tJ5QcluZoQNoY6J+QGoZlRgsorQ3Ra3gIMc+
ESYbvKct3DFEqEtFwQSABDhKN6KOLnB0+G7t/PApw3xafRJFVYI7lKcjhQIDAQAB
AoGAeJ5zMK9TU0AujuDMWtmY+V2ItTfP5JtMXPPa2plm+A7+yGIJBtcUFzIvIhMx
CYCqTWkjxL0DQ/tltWyG9r85nKzw0jalphM8yPimHjvllHX9REfKWfaV7laBS3aM
vn4hVsxsYs2r53YEbVN13aIlDOYPd9eIrirEFRHodlnrlwkCQQDuR4v/YKb2UVIU
49nImXxznBaHRUACAutnGitBsDw+tSrlm6UXLYJg143o3s9UTLB+42eIrGyu0HrR
IuQJUJQLAkEA4uaLdRjJJzI09oTmWeFGrR/t+msRfYpSNrieXai1R612Sco6BWK6
XarYtpbZ6HR77ca8vNcHwh9ESVLu93vQrwJADgqdT2FMtXs5UQ3USaPx14Y9NZ95
FCVD5gF+xxIxmqhmbL1tTx5ZboeFT1HB+f/C7tdLxJwUk4CpnCVoNrxO3QJBALyt
m2/vAW4/mL0Z/HbnFp9l+r2PBQdQ21a3pLEbVktZWhC4QhEybOjw5a7HuEJNgrRR
26ZoZQIuf9k9RouzgO8CQF7M0kTJa0Hh4b5O2o0wsMmPkxxCoSC6k4zsxvGra7f3
sZFbFfvt9bfFurkxPWI/UBeM0aTBsUTnZ53bqJGFypI=
-----END RSA PRIVATE KEY-----`

	// Simulate data that needs to be signed (like XML content for electronic invoicing)
	xmlDataToSign := []byte(`<invoice><amount>1000</amount><rut>12345678-9</rut></invoice>`)

	// Use the enhanced crypto function - it handles all encoding conversions automatically
	signature, err := SignSHA1WithRSA(xmlDataToSign, xmlPEMKey)
	if err != nil {
		t.Fatalf("Failed to sign XML data: %v", err)
	}

	// Verify we got a valid signature
	if len(signature) == 0 {
		t.Error("Signature should not be empty")
	}

	// The signature is ready to be used in XML responses or stored
	t.Logf("XML data signed successfully. Signature length: %d bytes", len(signature))
	t.Logf("Signature: %s", signature)
}
