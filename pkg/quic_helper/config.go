package quic_helper

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"

	"github.com/lucas-clemente/quic-go"
)

// Setup a bare-bones TLS config for the server
func GetGenerateTLSConfig(token string) *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates:           []tls.Certificate{tlsCert},
		NextProtos:             []string{token},
		SessionTicketsDisabled: true,
	}
}
func GetClientTlsConfig(otherToken string) *tls.Config {
	return &tls.Config{
		InsecureSkipVerify:     true,
		NextProtos:             []string{otherToken},
		SessionTicketsDisabled: true,
	}
}

func GetQUICConfig() *quic.Config {
	return &quic.Config{
		ConnectionIDLength: 4,
		KeepAlive:          true,
	}
}
