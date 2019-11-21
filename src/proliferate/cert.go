package proliferate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"path"
	"time"
)

func (node *Node) GenerateX509Pair() (string, string) {
	n := *node

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	}

	pemKey := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	tml := x509.Certificate{
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(5, 0, 0),

		SerialNumber: big.NewInt(IssueSerial()),
		Subject: pkix.Name{
			CommonName:   n.Detail.Name,
			Organization: []string{n.Detail.Organization},
		},
		BasicConstraintsValid: true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, &tml, &tml, &key.PublicKey, key)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	}

	pemCert := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	return string(pemKey), string(pemCert)
}

// TODO needs to return actual serial
func IssueSerial() int64 {
	return 1
}

func (node *Node) IdentityCertificates() (string, string) {
	n := *node
	c := n.Config.Static

	certFile := path.Join(c.IdentityFolder, c.CertFile)
	keyFile := path.Join(c.IdentityFolder, c.KeyFile)

	return certFile, keyFile
}

// CertificateLoad attaches node certificates to n.member
func (node *Node) IdentityCertificateLoad() {
	n := *node

	certFile, keyFile := n.IdentityCertificates()

	cert, err := ioutil.ReadFile(certFile)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	} else {
		n.member.cert = string(cert)
	}

	key, err := ioutil.ReadFile(keyFile)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	} else {
		n.member.key = string(key)
	}

	*node = n
}
