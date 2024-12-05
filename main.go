package main

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"github.com/thibauult/tee-mock-server/pki"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	socketPath = "/run/container_launcher/teeserver.sock"
	// default values
	defaultGoogleServiceAccount     = "tee-mock-server@localhost.gserviceaccount.com"
	defaultTokenExpirationInMinutes = 5
)

type attestationTokenRequest struct {
	Audience  string   `json:"audience"`
	TokenType string   `json:"token_type"`
	Nonces    []string `json:"nonces"`
}

type tokenConfig struct {
	signingKey               *rsa.PrivateKey
	chain                    interface{}
	googleServiceAccount     string
	tokenExpirationInMinutes int
}

func main() {

	config := loadConfig()

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
	config.signingKey, err = pki.GetSigningPrivateKey()
	if err != nil {
		log.Printf("Failed to load private key: %v\n", err)
		log.Fatal(err)
	}

	config.chain = pki.GetCertificateChain()

	log.Printf("Root Certificate:\n\n%s\n", pki.GetRootCertificate())

	m := http.NewServeMux()
	m.HandleFunc("/v1/token", newPostTokenHandler(config))

	server := http.Server{Handler: m}

	log.Println("Start serving on socket", socketPath, "...")
	if err := server.Serve(socket); err != nil {
		log.Fatal(err)
	}
}

func loadConfig() tokenConfig {

	const googleServiceAccount = "google_service_account"
	const tokenExpirationInMinutes = "token_expiration_in_minutes"

	viper.SetDefault(googleServiceAccount, defaultGoogleServiceAccount)
	viper.SetDefault(tokenExpirationInMinutes, defaultTokenExpirationInMinutes)

	viper.SetEnvPrefix("tee")
	err := viper.BindEnv(tokenExpirationInMinutes, googleServiceAccount)
	if err != nil {
		log.Fatal(err)
	}

	config := tokenConfig{}
	config.googleServiceAccount = viper.GetString(googleServiceAccount)
	config.tokenExpirationInMinutes = viper.GetInt(tokenExpirationInMinutes)

	log.Println(config)

	return config
}

func newPostTokenHandler(config tokenConfig) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, _ = fmt.Fprintf(w, "{ \"message\": \"only POST is supported\" }")
			return
		}

		tokenRequest, err := parseAndValidateAttestationTokenRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "{ \"message\": \"%s\" }", err.Error())
			return
		}

		signedJwt := newToken(config, tokenRequest)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(signedJwt + "\n"))
	}
}

func parseAndValidateAttestationTokenRequest(r *http.Request) (*attestationTokenRequest, error) {

	tokenRequest := attestationTokenRequest{}
	err := json.NewDecoder(r.Body).Decode(&tokenRequest)
	if err != nil {
		return nil, errors.New("failed to parse attestation token request")
	}

	if tokenRequest.TokenType != "PKI" {
		return nil, fmt.Errorf("invalid token type: %s", tokenRequest.TokenType)
	}

	if len(tokenRequest.Audience) == 0 {
		return nil, fmt.Errorf("audience not set")
	}

	return &tokenRequest, nil
}

func newToken(config tokenConfig, tokenRequest *attestationTokenRequest) string {
	log.Println("Creating new Token...")

	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, createClaims(config, tokenRequest.Audience, tokenRequest.Nonces))
	token.Header["x5c"] = config.chain

	// Sign the token with the RSA private key
	signedToken, err := token.SignedString(config.signingKey)
	if err != nil {
		log.Printf("Error signing token: %v\n", err)
	}

	log.Println("New Token successfully signed")
	return signedToken
}

func createClaims(config tokenConfig, audience string, nonces []string) jwt.MapClaims {

	if nonces == nil {
		nonces = []string{}
	}

	return jwt.MapClaims{
		"iss":       "https://confidentialcomputing.googleapis.com",
		"sub":       "https://www.googleapis.com/compute/v1/projects/PROJECT_ID/zones/us-central1-a/instances/INSTANCE_NAME", // TODO set PROJECT_ID and INSTANCE_NAME
		"aud":       audience,
		"eat_nonce": nonces,
		"exp":       time.Now().Add(5 * time.Minute).Unix(),
		"iat":       time.Now().Unix(),
		"nbf":       time.Now().Unix(),
		// Confidential Space Claims
		"eat_profile": "https://cloud.google.com/confidential-computing/confidential-space/docs/reference/token-claims",
		"secboot":     true,
		"oemid":       11129,
		"google_service_accounts": []string{
			config.googleServiceAccount,
		},
		"hwmodel":   "GCP_AMD_SEV",
		"swname":    "CONFIDENTIAL_SPACE",
		"swversion": []string{"240900"},
		"submods": map[string]interface{}{
			"confidental_space": map[string]interface{}{
				"monitoring_enabled": map[string]bool{
					"memory": false,
				},
				"support_attributes": []string{
					"LATEST", "STABLE", "USABLE",
				},
			},
			"container": map[string]interface{}{
				"args": []string{
					"/customnonce",
					"/docker-entrypoint.sh",
					"nginx",
					"-g",
					"daemon off;",
				},
				"env": map[string]string{
					"HOSTNAME":      "HOST_NAME", // TODO set HOST_NAME
					"NGINX_VERSION": "1.27.0",
					"NJS_RELEASE":   "2~bookworm",
					"NJS_VERSION":   "0.8.4",
					"PATH":          "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
					"PKG_RELEASE":   "2~bookworm",
				},
				"image_digest":    "sha256:67682bda769fae1ccf5183192b8daf37b64cae99c6c3302650f6f8bf5f0f95df",
				"image_id":        "sha256:fffffc90d343cbcb01a5032edac86db5998c536cd0a366514121a45c6723765c",
				"image_reference": "docker.io/library/nginx:latest",
				"image_signatures": []interface{}{
					map[string]string{
						"key_id":              "<hexadecimal-sha256-fingerprint-public-key1>",
						"signature":           "<base64-encoded-signature>",
						"signature_algorithm": "RSASSA_PSS_SHA256",
					},
					map[string]string{
						"key_id":              "<hexadecimal-sha256-fingerprint-public-key2>",
						"signature":           "<base64-encoded-signature>",
						"signature_algorithm": "RSASSA_PSS_SHA256",
					},
					map[string]string{
						"key_id":              "<hexadecimal-sha256-fingerprint-public-key3>",
						"signature":           "<base64-encoded-signature>",
						"signature_algorithm": "ECDSA_P256_SHA256",
					},
				},
				"restart_policy": "Never",
			},
			"gce": map[string]string{
				"instance_id":    "INSTANCE_ID",    // TODO set INSTANCE_ID
				"instance_name":  "INSTANCE_NAME",  // TODO set INSTANCE_NAME
				"project_id":     "PROJECT_ID",     // TODO set PROJECT_ID
				"project_number": "PROJECT_NUMBER", // TODO set PROJECT_NUMBER
				"zone":           "us-central1-a",  // TODO set ZONE
			},
		},
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
