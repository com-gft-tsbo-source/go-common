package dispatcher

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
)

// TLSInfo ...
type TLSInfo struct {
	certificate *tls.Certificate
	caCertPool  *x509.CertPool
	tlsConfig   *tls.Config
	transport   *http.Transport
}

func checkCertificate(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {

	var msg string
	var hosts string
	var ips string

	for _, rawCert := range rawCerts {

		cert, _ := x509.ParseCertificate(rawCert)

		for _, host := range cert.DNSNames {
			if len(host) == 0 {
				continue
			}
			hosts = fmt.Sprintf("%s %s", hosts, host)
		}

		for _, ip := range cert.IPAddresses {
			if len(ip) == 0 {
				continue
			}
			ips = fmt.Sprintf("%s %s", ips, ip)
		}

	}

	if len(hosts) > 0 {
		hosts = fmt.Sprintf("hosts [%s]", hosts[1:])
		msg = fmt.Sprintf("Found certificates for %s", hosts)
	}

	if len(ips) > 0 {
		ips = fmt.Sprintf("ips [%s]", ips[1:])
		if len(msg) == 0 {
			msg = fmt.Sprintf("Found certificates for %s", ips)
		} else {
			msg = fmt.Sprintf("%s and %s", msg, ips)
		}
	}
	fmt.Println(msg)
	return nil
}
