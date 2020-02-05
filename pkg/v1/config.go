package v1

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// IAM ..
type IAM struct {
	ClientSecret string `yaml:"client_secret"`
	ClientID     string `yaml:"client_id"`
	Token        string `yaml:"token"`
	Endpoint     string `yaml:"endpoint"`
	Credentials  string `yaml:"credentials"`
}

// TTSClient ..
type TTSClient struct {
	IAMClient IAM `yaml:"iam"`
}

// GetConf ..
func (c *TTSClient) GetConf(path string) (*TTSClient, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return &TTSClient{}, err
	}

	log.Printf("Reading: %s", path)
	err = yaml.UnmarshalStrict(f, &c)
	if err != nil {
		return &TTSClient{}, err
	}

	return c, nil
}

// CreateCertSecret ..
func (c *TTSClient) CreateCertSecret(id string, entries []ProxyEntries, kubeClientset *kubernetes.Clientset) (*v1.Secret, error) {
	// if secret exists remove it
	_, err := kubeClientset.CoreV1().Secrets(v1.NamespaceDefault).Get("certs-secret", metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return &v1.Secret{}, fmt.Errorf("Unexpected problem checking of secret existance: %s", err)
	} else if err == nil {
		log.Print("Cert secret already exists, removing it before proceeding")
		err = kubeClientset.CoreV1().Secrets(v1.NamespaceDefault).Delete("certs-secret", &metav1.DeleteOptions{})
		if err != nil {
			return &v1.Secret{}, fmt.Errorf("Failed to delete previous secret: %s", err)
		}
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "certs-secret",
			Namespace: v1.NamespaceDefault,
		},
		Data: map[string][]byte{
			//"watts_cert.key": []byte(refreshToken),
		},
	}

	for _, entry := range entries {
		name := strings.ReplaceAll(entry.Name, " ", "")
		name = strings.Split(name, "(")[0]
		secret.Data[name] = []byte(entry.Value)
	}

	secret.Data["id"] = []byte(id)

	certSecret, err := kubeClientset.CoreV1().Secrets(v1.NamespaceDefault).Create(secret)
	if err != nil {
		return &v1.Secret{}, fmt.Errorf("Failed to create cert secret: %s", err)
	}

	return certSecret, nil
}

// CreateSecret ..
func (c *TTSClient) CreateSecret(refreshToken string, kubeClientset *kubernetes.Clientset) (*v1.Secret, error) {

	// if secret exists remove it
	_, err := kubeClientset.CoreV1().Secrets(v1.NamespaceDefault).Get("refresh-token", metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return &v1.Secret{}, fmt.Errorf("Unexpected problem checking of secret existance: %s", err)
	} else if err == nil {
		log.Print("Token secret already exists, removing it before proceeding")
		err = kubeClientset.CoreV1().Secrets(v1.NamespaceDefault).Delete("refresh-token", &metav1.DeleteOptions{})
		if err != nil {
			return &v1.Secret{}, fmt.Errorf("Failed to delete previous secret: %s", err)
		}
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "refresh-token",
			Namespace: v1.NamespaceDefault,
		},
		Data: map[string][]byte{
			"RefreshToken": []byte(refreshToken),
			"AccessToken":  []byte(c.IAMClient.Token),
		},
	}

	tokenSecret, err := kubeClientset.CoreV1().Secrets(v1.NamespaceDefault).Create(secret)
	if err != nil {
		return &v1.Secret{}, fmt.Errorf("Failed to create token secret: %s", err)
	}

	log.Printf("Refresh token in secret: %s", tokenSecret.Data["RefreshToken"])

	return tokenSecret, nil

}
