package certutil

import (
	"crypto/x509"
	"math/big"
	"time"
)

// CertOption is a modification of a certificate.
type CertOption func(*x509.Certificate) error

// OptSerialNumber sets the certificate serial number.
func OptSerialNumber(serialNumber *big.Int) CertOption {
	return func(csr *x509.Certificate) error {
		csr.SerialNumber = serialNumber
		return nil
	}
}

// OptSubjectCommonName sets the subject common name.
func OptSubjectCommonName(commonName string) CertOption {
	return func(csr *x509.Certificate) error {
		csr.Subject.CommonName = commonName
		return nil
	}
}

// OptSubjectOrganization sets the subject organization names.
func OptSubjectOrganization(organization ...string) CertOption {
	return func(csr *x509.Certificate) error {
		csr.Subject.Organization = organization
		return nil
	}
}

// OptSubjectCountry sets the subject country names.
func OptSubjectCountry(country ...string) CertOption {
	return func(csr *x509.Certificate) error {
		csr.Subject.Country = country
		return nil
	}
}

// OptSubjectProvince sets the subject province names.
func OptSubjectProvince(province ...string) CertOption {
	return func(csr *x509.Certificate) error {
		csr.Subject.Province = province
		return nil
	}
}

// OptSubjectLocality sets the subject locality names.
func OptSubjectLocality(locality ...string) CertOption {
	return func(csr *x509.Certificate) error {
		csr.Subject.Locality = locality
		return nil
	}
}

// OptNotAfter sets the not after time.
func OptNotAfter(notAfter time.Time) CertOption {
	return func(csr *x509.Certificate) error {
		csr.NotAfter = notAfter
		return nil
	}
}

// OptNotBefore sets the not before time.
func OptNotBefore(notBefore time.Time) CertOption {
	return func(csr *x509.Certificate) error {
		csr.NotBefore = notBefore
		return nil
	}
}

// OptIsCA sets the is certificate authority flag.
func OptIsCA(isCA bool) CertOption {
	return func(csr *x509.Certificate) error {
		csr.IsCA = isCA
		return nil
	}
}

// OptKeyUsage sets the key usage flags.
func OptKeyUsage(keyUsage x509.KeyUsage) CertOption {
	return func(csr *x509.Certificate) error {
		csr.KeyUsage = keyUsage
		return nil
	}
}

// OptDNSNames sets valid dns names for the cert.
func OptDNSNames(dnsNames ...string) CertOption {
	return func(csr *x509.Certificate) error {
		csr.DNSNames = dnsNames
		return nil
	}
}

// OptAdditionalNames adds valid dns names for the cert.
func OptAdditionalNames(dnsNames ...string) CertOption {
	return func(csr *x509.Certificate) error {
		csr.DNSNames = append(csr.DNSNames, dnsNames...)
		return nil
	}
}
