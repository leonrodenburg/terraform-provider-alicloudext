package certificates

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"strings"

	"github.com/go-acme/lego/v3/certcrypto"
	"github.com/go-acme/lego/v3/certificate"
	"github.com/go-acme/lego/v3/challenge"
	"github.com/go-acme/lego/v3/lego"
	"github.com/go-acme/lego/v3/providers/dns/alidns"
	"github.com/go-acme/lego/v3/registration"
)

type User struct {
	Email        string
	registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u User) GetRegistration() *registration.Resource {
	return u.registration
}

func (u *User) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func RequestCertificateForDomainUsingDns(domain string, client *lego.Client, provider challenge.Provider) (*certificate.Resource, error) {
	certRequest := certificate.ObtainRequest{
		Domains: []string{domain},
		Bundle:  true,
	}

	err := client.Challenge.SetDNS01Provider(provider)
	if err != nil {
		return nil, err
	}

	issued, err := client.Certificate.Obtain(certRequest)
	if err != nil {
		return nil, err
	}

	return issued, nil
}

func CreateAlicloudDnsProvider(accessKey string, secretKey string, region string) (*alidns.DNSProvider, error) {
	alidnsConfig := alidns.NewDefaultConfig()
	alidnsConfig.APIKey = accessKey
	alidnsConfig.SecretKey = secretKey
	alidnsConfig.RegionID = region
	return alidns.NewDNSProviderConfig(alidnsConfig)
}

func CreateCertClientForUser(user *User) (*lego.Client, error) {
	privateKey, err := getNewPrivateKey()
	if err != nil {
		return nil, err
	}
	user.key = privateKey
	config := lego.NewConfig(user)
	config.CADirURL = "https://acme-v02.api.letsencrypt.org/directory"
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		return nil, err
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	if err != nil {
		return nil, err
	}
	user.registration = reg

	return client, nil
}

func getNewPrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func Sanitize(s string) string {
	return strings.Replace(s, "\n\n", "\n", -1)
}