package main

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"github.com/thibauult/tee-mock-server/pki"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const socketPath = "/run/container_launcher/teeserver.sock"

func main() {

	// Create a Unix domain socket and listen for incoming connections.
	socket, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal(err)
	}

	// Cleanup the socket file.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go removeSocketOnInterrupt(c)

	// Load RSA private key
	signingKey, err := pki.GetSigningPrivateKey()
	if err != nil {
		log.Printf("Failed to load private key: %v\n", err)
		log.Fatal(err)
	}

	chain := pki.GetCertificateChain()

	log.Printf("Root Certificate:\n\n%s\n", pki.GetRootCertificate())

	m := http.NewServeMux()
	m.HandleFunc("/v1/token", func(w http.ResponseWriter, r *http.Request) {
		signedJwt := newToken(signingKey, chain)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(signedJwt + "\n"))
	})

	server := http.Server{Handler: m}

	log.Println("Start serving on socket", socketPath, "...")
	if err := server.Serve(socket); err != nil {
		log.Fatal(err)
	}
}

func newToken(signingKey *rsa.PrivateKey, chain []string) string {
	log.Println("Creating new Token...")

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, createClaims())
	token.Header["x5c"] = chain

	// Sign the token with the RSA private key
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		log.Printf("Error signing token: %v\n", err)
	}

	log.Println("New Token successfully signed")
	return signedToken
}

func createClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"iss": "https://confidentialcomputing.googleapis.com",
		"sub": "https://www.googleapis.com/compute/v1/projects/PROJECT_ID/zones/us-central1-a/instances/INSTANCE_NAME",
		"aud": "AUDIENCE_NAME",
		"exp": time.Now().Add(5 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		// Confidential Space Claims
		"eat_profile": "https://cloud.google.com/confidential-computing/confidential-space/docs/reference/token-claims",
		"secboot":     true,
		"oemid":       11129,
		"google_service_accounts": []string{
			"tee-mock-server@localhost.gserviceaccount.com",
		},
		"hwmodel":   "GCP_AMD_SEV",
		"swname":    "CONFIDENTIAL_SPACE",
		"swversion": []string{"240900"},
		"submods":   map[string]interface{}{},
	}
}

func removeSocketOnInterrupt(c chan os.Signal) {
	<-c
	log.Println("Removing socket...")
	err := os.Remove(socketPath)
	if err != nil {
		log.Println("Failed to remove", socketPath, err)
	} else {
		log.Println("Successfully removed", socketPath)
	}
	os.Exit(1)
}
