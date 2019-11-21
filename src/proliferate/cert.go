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

// GenerateX509Pair returns X.509 key and certificate
func (node *Node) GenerateX509Pair() (string, string) {
	n := *node
	y := n.Config.Build.CertExpYears

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

	template := x509.Certificate{
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(y, 0, 0),

		SerialNumber: big.NewInt(IssueSerial()),
		Subject: pkix.Name{
			CommonName:   n.Detail.Name,
			Organization: []string{n.Detail.Organization},
		},
		BasicConstraintsValid: true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
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

// IdentityCertificates returns filesystem paths to identy cert and key
func (node *Node) IdentityCertificates() (string, string) {
	n := *node
	c := n.Config.Build

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

func (node *Node) VerifySignature(certPEM string, rootPEM string) {
	n := *node

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		n.Log(Message{
			Level: 2,
			Text:  "Failed to parse root certificate",
		})
	}

	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		n.Log(Message{
			Level: 2,
			Text:  "failed to parse certificate PEM",
		})
	}
	cert, err := x509.ParseCertificate(block.Bytes)

	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	}

	opts := x509.VerifyOptions{
		//DNSName: "www.domain.com",
		Roots: roots,
	}

	if _, err := cert.Verify(opts); err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	}
}
