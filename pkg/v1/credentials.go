package v1

import (
	"encoding/json"
	"fmt"
	"log"

	iam "github.com/dodas-ts/dodas-go-client"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
)

// ProxyEntries ..
type ProxyEntries struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	SaveAs string `json:"save_as,omitempty"`
	Type   string `json:"type"`
	Rows   int    `json:"rows,omitempty"`
	Cols   int    `json:"cols,omitempty"`
}

// CredEntries ..
type CredEntries struct {
	Entries []ProxyEntries `json:"entries"`
	ID      string         `json:"id"`
}

// ProxyStruct ..
type ProxyStruct struct {
	Credentials CredEntries `json:"credential"`
	Result      string      `json:"result"`
}

// CacheCredentials ..
func (c *TTSClient) CacheCredentials(kubeClientset *kubernetes.Clientset) error {
	iamClient := iam.Conf{
		AllowRefresh: iam.TokenRefreshConf{
			ClientID:         ttsClient.IAMClient.ClientID,
			ClientSecret:     ttsClient.IAMClient.ClientSecret,
			IAMTokenEndpoint: ttsClient.IAMClient.Endpoint,
		},
		Im: iam.ConfIM{
			Token: ttsClient.IAMClient.Token,
		},
	}

	// if secret exists skip refresh token
	tokenSecret, err := kubeClientset.CoreV1().Secrets(v1.NamespaceDefault).Get("refresh-token", metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("Unexpected problem checking of secret existance: %s", err)
	} else if err == nil {
		log.Print("Token secret already exists, skipping refresh token retrieval")
	} else {

		refreshToken, err := iamClient.GetRefreshToken()
		if err != nil {
			return fmt.Errorf("Failed to retrieve RefreshToken: %s", err)
		}

		tokenSecret, err = ttsClient.CreateSecret(refreshToken)
		if err != nil {
			return fmt.Errorf("Failed to create Token secret: %s", err)
		}

		log.Printf("Refresh token in secret: %s", tokenSecret.Data["RefreshToken"])
	}

	// TODO: check it token is valid before refreshing
	//accessToken := string(tokenSecret.Data["AccessToken"])

	//log.Printf("Refreshing access token \n %s \n %s", tokenSecret.Data["RefreshToken"], tokenSecret.Data["AccessToken"])

	accessToken, err := iamClient.GetAccessToken(string(tokenSecret.Data["RefreshToken"]))
	if err != nil {
		return fmt.Errorf("Failed get Refresh AccessToken: %s", err)
	}

	//log.Printf("New access token: %s", accessToken)

	log.Printf("Updating new access token in secret")

	tokenSecret.Data["AccessToken"] = []byte(accessToken)
	tokenSecret, err = kubeClientset.CoreV1().Secrets(v1.NamespaceDefault).Update(tokenSecret)
	if err != nil {
		return fmt.Errorln(err)
	}

	//log.Printf("Refreshing access token \n %s \n %s", tokenSecret.Data["RefreshToken"], tokenSecret.Data["AccessToken"])

	log.Print("New access token saved.")

	// retrieve credential info and get proxy
	// TODO: if proxy is valid skip
	// TODO: revoke previous one

	request := iam.Request{
		URL:         ttsClient.IAMClient.Credentials,
		RequestType: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + accessToken,
			"Content-Type":  "application/json",
		},
		Content: []byte("{\"service_id\": \"x509\"}"),
	}

	body, statusCode, err := iam.MakeRequest(request)
	if err != nil {
		return fmt.Errorln(err)
	}

	if statusCode == 200 {
		proxyEntry := ProxyStruct{}

		err = json.Unmarshal(body, &proxyEntry)
		if err != nil {
			fmt.Errorln(err)
		}

		for _, entry := range proxyEntry.Credentials.Entries {
			log.Println(entry.Name)
			log.Println(entry.Value)
		}

		_, err = ttsClient.CreateCertSecret(proxyEntry.Credentials.Entries)
		if err != nil {
			fmt.Errorln(err)
		}

	} else {
		return fmt.Errorf("code %d: %s", statusCode, body)
	}

	return nil
	// TODO: Set cert and key in k8s certs secret
}
