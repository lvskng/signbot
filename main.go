package main

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type SignRequest struct {
	Message string `json:"message"` //ASCII encoded param string
	Key     string `json:"key"`     //in base64
}

type VerifyRequest struct {
	Message   string `json:"message"`
	Signature string `json:"signature"`
	Key       string `json:"key"`
}

type SignResponse struct {
	Signature string `json:"signature"` //in base64
}

type VerifyResponse struct {
	Valid bool `json:"valid"`
}

var (
	port = flag.Uint("port", 8080, "port to listen on")
)

func main() {
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, err := fmt.Fprintf(w, "alive")
		if err != nil {
			return
		}
	})
	http.HandleFunc("/sign", sign)
	http.HandleFunc("/verify", verify)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		log.Fatal(err)
	}
}

func sign(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req SignRequest
	err := decoder.Decode(&req)
	if err != nil {
		if rs, e := io.ReadAll(r.Body); e == nil {
			log.Printf("error unmarshaling request %v: %v", rs, err.Error())
		} else {
			log.Printf("error unmarshaling request: %v", err.Error())
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	key, err := getKeyFromString(req.Key)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "Error decoding key", http.StatusBadRequest)
		return
	}

	signed := ed25519.Sign(key, []byte(req.Message))
	resp := SignResponse{
		Signature: base64.StdEncoding.EncodeToString(signed),
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func getKeyFromString(keyString string) (ed25519.PrivateKey, error) {
	if !strings.HasPrefix(keyString, "-----BEGIN PRIVATE KEY-----") {
		keyString = "-----BEGIN PRIVATE KEY-----\n" + keyString
	}
	if !strings.HasSuffix(keyString, "-----END PRIVATE KEY-----") {
		keyString += "\n-----END PRIVATE KEY-----"
	}
	block, _ := pem.Decode([]byte(keyString))
	if block == nil {
		return nil, fmt.Errorf("error decoding key: invalid PEM block")
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil || reflect.ValueOf(parsed).IsNil() {
		return nil, fmt.Errorf("error decoding key: %v", err)
	} else if reflect.ValueOf(parsed).Type() != reflect.TypeOf(ed25519.PrivateKey{}) {
		return nil, fmt.Errorf("error decoding key: expected ed25519.PrivateKey, got %T", parsed)
	}
	key := parsed.(ed25519.PrivateKey)
	if key == nil {
		return nil, fmt.Errorf("error decoding key: expected ed25519.PrivateKey, got %T", parsed)
	}
	return key, nil
}

func verify(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var req VerifyRequest
	err := decoder.Decode(&req)
	if err != nil {
		if rs, e := io.ReadAll(r.Body); e == nil {
			log.Printf("error unmarshaling request %v: %v", rs, err.Error())
		} else {
			log.Printf("error unmarshaling request: %v", err.Error())
		}
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	key, err := getKeyFromString(req.Key)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "Error decoding key", http.StatusBadRequest)
		return
	}
	pubKey := key.Public().(ed25519.PublicKey)
	if pubKey == nil {
		log.Printf("error decoding key: expected ed25519.PublicKey, got %T", pubKey)
		http.Error(w, "Error decoding key", http.StatusBadRequest)
		return
	}
	sig, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		log.Printf("error decoding message: %v", err)
		http.Error(w, "Message not in Base64 format", http.StatusBadRequest)
		return
	}
	resp := VerifyResponse{Valid: ed25519.Verify(pubKey, []byte(req.Message), sig)}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
