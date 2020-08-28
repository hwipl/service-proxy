package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestParseTCPAddr(t *testing.T) {
	var want, got []string

	addrs := []string{
		"127.0.0.1",
		"192.168.1.1:3232",
		"::1",
		"[2000::1]:2323",
	}
	want = []string{
		"127.0.0.1:32323",
		"192.168.1.1:3232",
		"[::1]:32323",
		"[2000::1]:2323",
	}
	for _, i := range addrs {
		got = append(got, parseTCPAddr(i).String())
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func createTestFile(name string, content []byte) *os.File {
	tmpFile, err := ioutil.TempFile("", name)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := tmpFile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		log.Fatal(err)
	}
	return tmpFile
}

func createTestCertFile(name string) *os.File {
	certPem := []byte(`-----BEGIN CERTIFICATE-----
MIIBhTCCASugAwIBAgIQIRi6zePL6mKjOipn+dNuaTAKBggqhkjOPQQDAjASMRAw
DgYDVQQKEwdBY21lIENvMB4XDTE3MTAyMDE5NDMwNloXDTE4MTAyMDE5NDMwNlow
EjEQMA4GA1UEChMHQWNtZSBDbzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABD0d
7VNhbWvZLWPuj/RtHFjvtJBEwOkhbN/BnnE8rnZR8+sbwnc/KhCk3FhnpHZnQz7B
5aETbbIgmuvewdjvSBSjYzBhMA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAKBggr
BgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MCkGA1UdEQQiMCCCDmxvY2FsaG9zdDo1
NDUzgg4xMjcuMC4wLjE6NTQ1MzAKBggqhkjOPQQDAgNIADBFAiEA2zpJEPQyz6/l
Wf86aX6PepsntZv2GYlA5UpabfT2EZICICpJ5h/iI+i341gBmLiAFQOyTDT+/wQc
6MF9+Yw1Yy0t
-----END CERTIFICATE-----`) // from go tls examples
	return createTestFile(name, certPem)
}

func createTestKeyFile(name string) *os.File {
	keyPem := []byte(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIIrYSSNQFaA2Hwf1duRSxKtLYX5CB04fSeQ6tF1aY/PuoAoGCCqGSM49
AwEHoUQDQgAEPR3tU2Fta9ktY+6P9G0cWO+0kETA6SFs38GecTyudlHz6xvCdz8q
EKTcWGekdmdDPsHloRNtsiCa697B2O9IFA==
-----END EC PRIVATE KEY-----`) // from go tls examples
	return createTestFile(name, keyPem)
}

func TestParseCertFiles(t *testing.T) {
	// create certificate file
	cf := createTestCertFile("parsecertfilestest-cert-*.pem")
	defer os.Remove(cf.Name())

	// create key file
	kf := createTestKeyFile("parsecertfilestest-key-*.pem")
	defer os.Remove(kf.Name())

	// test parsing
	certFile = cf.Name()
	keyFile = kf.Name()
	parseCertFiles()
}

func TestParseCACertFiles(t *testing.T) {
	// create certificate file
	cf := createTestCertFile("parsecacertfilestest-cert-*.pem")
	defer os.Remove(cf.Name())

	// test single file parsing
	caCertFiles = cf.Name()
	parseCACertFiles()

	// test multiple files parsing
	caCertFiles += "," + cf.Name() + "," + cf.Name()
	parseCACertFiles()
}
