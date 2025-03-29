package main

func main() {
	// Basic Options
	DefiningOptionalValues()
	InspectingValues()
	MarshalingExamples()
	TransformationExamples()

	// File Options
	path := "examples/readme_examples/example.txt"
	LoadingAndReadingFiles(&path)
	WritingAndDeletingFiles(&path)
	AdditionalFileTools(&path)

	// Cert and Key Options
	certPath := "testing/rsa/cert.pem"
	keyPath := "testing/rsa/key.pem"
	pubKeyPath := "testing/ed25519/pub.pem"
	privKeyPath := "testing/ed25519/key.pem"
	TLSCertificate(&certPath, &keyPath)
	SigningKeys(&pubKeyPath, &privKeyPath)
}
