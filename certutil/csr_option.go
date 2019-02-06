package certutil

import (
	"crypto/x509"
	"time"
)

// CertOption is a modification of a certificate.
type CertOption func(csr *x509.Certificate)

// OptSubjectCommonName sets the subject common name.
func OptSubjectCommonName(commonName string) CertOption {
	return func(csr *x509.Certificate) {
		csr.Subject.CommonName = commonName
	}
}

// OptSubjectOrganization sets the subject organization names.
func OptSubjectOrganization(organization ...string) CertOption {
	return func(csr *x509.Certificate) {
		csr.Subject.Organization = organization
	}
}

// OptSubjectCountry sets the subject country names.
func OptSubjectCountry(country ...string) CertOption {
	return func(csr *x509.Certificate) {
		csr.Subject.Country = country
	}
}

// OptSubjectProvince sets the subject province names.
func OptSubjectProvince(province ...string) CertOption {
	return func(csr *x509.Certificate) {
		csr.Subject.Province = province
	}
}

// OptSubjectLocality sets the subject locality names.
func OptSubjectLocality(locality ...string) CertOption {
	return func(csr *x509.Certificate) {
		csr.Subject.Locality = locality
	}
}

// OptNotAfter sets the not after time.
func OptNotAfter(notAfter time.Time) CertOption {
	return func(csr *x509.Certificate) {
		csr.NotAfter = notAfter
	}
}

// OptNotBefore sets the not before time.
func OptNotBefore(notBefore time.Time) CertOption {
	return func(csr *x509.Certificate) {
		csr.NotBefore = notBefore
	}
}

// OptIsCA sets the is certificate authority flag.
func OptIsCA(isCA bool) CertOption {
	return func(csr *x509.Certificate) {
		csr.IsCA = isCA
	}
}

// OptKeyUsage sets the key usage flags.
func OptKeyUsage(keyUsage x509.KeyUsage) CertOption {
	return func(csr *x509.Certificate) {
		csr.KeyUsage = keyUsage
	}
}

// OptDNSNames sets valid dns names for the cert.
func OptDNSNames(dnsNames ...string) CertOption {
	return func(csr *x509.Certificate) {
		csr.DNSNames = dnsNames
	}
}

// OptAdditionalNames adds valid dns names for the cert.
func OptAdditionalNames(dnsNames ...string) CertOption {
	return func(csr *x509.Certificate) {
		csr.DNSNames = append(csr.DNSNames, dnsNames...)
	}
}
