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
	// Generated with: openssl genrsa 512
	testPrivateKey := `-----BEGIN RSA PRIVATE KEY-----
MIIBOgIBAAJBAMyuWln6lOmrlwCB0bHjOeeO0hjFWJ3qaNlkklqc0PBFkphymxfk
AwoG0Cdra977wRASlZut9tjCDN03a2qvclMCAwEAAQJAP7SsroDNxIUBtMizKbjd
lvLe9ZLG6C/DfpZM7yML7RxL6U3+zZ+OVGE0QtJ1Q5PN5kYS9HuFr4yAEU5HTPNb
kQIhAOertvNvydi0+ZdOMi72HvbxrHbw8Ej4eKp8o0vskwgVAiEA4i0Kq4B2yWuF
vyLIvIUmXZu53V+iXiS7vzn4PcYEAscCIQComn/7i1ALNyquw2oiY10Fu70YkyFM
+ghXi34Ms5AOQQIgXpogYGO3S8BhjPTrqY634WeFcobRzzbmILIKlyv/+XkCIG0m
sbFMkB/JdjojoZ6FSFiUpW+ItmRuIr4DXTlv2YS7
-----END RSA PRIVATE KEY-----`

	// 	testPrivateKey := `-----BEGIN RSA PRIVATE KEY-----
	// MIIBOgIBAAJBAMY4zzWGBVD77NaFWei0U7VPmMGcShFVwj1KFTgzftRKxeTrqBhr
	// +vuU0DyZULm2NLV9ERcCn6mTxU/sFdO17KkCAQMCQQC9KqZzUFz+0n+qT+4+J1Cb
	// +kT+UxZ6uR/yt2c/Af03Gy7hWXWI7W3h3FLWJCb+8+JX1xZ1qH2W9WXH7bqPvWQb
	// AgEAAoGAJHoWH8QP1fM8hWE5DcS0Qy2e2XUQZhN5QZhV0oqEQs2jc+rpqaHYrGVB
	// XZZ+3AJZx5BwkFUaEmpxQX1YZWKFnNK6Nh7G4y7ZI4Yv7S5/3RxXDG7OUvGtzY7O
	// y5Qr6LwjCtZBQJ5Ef5MxhXv8HhoPOrqoHZZGmZi/1+0FvpYDZQECQQDGOc81hgVQ
	// ++zWhVnotFO1T5jBnEoRVcI9ShU4M37USsXk66gYa/r7lNA8mVC5tjS1fREXAp+p
	// k8VP7BXTteypAkEAvSqmc1Bc/tJ/qk/uPidQm/pE/lMWerkf8rdnPwH9Nxsu4Vl1
	// iO1t4dxS1iQm/vPiV9cWdah9lvVlx+26j71kGw==
	// -----END RSA PRIVATE KEY-----`

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
