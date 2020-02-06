package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	iam "github.com/dodas-ts/dodas-go-client/cmd"
	"k8s.io/client-go/kubernetes"
)

// MapUser ..
func (c *TTSClient) MapUser(user string, accessToken string, kubeClientset *kubernetes.Clientset) error {

	request := iam.Request{
		URL:         c.IAMClient.Credentials,
		RequestType: "POST",
		Headers: map[string]string{
			"Authorization": "Bearer " + accessToken,
			"Content-Type":  "application/json",
		},
		Content: []byte("{\"service_id\": \"x509\"}"),
	}

	body, statusCode, err := iam.MakeRequest(request)
	if err != nil {
		return fmt.Errorf("Error getting user credentials: %s", err)
	}

	var secrets = make(map[string][]byte)

	if statusCode == 200 {
		proxyEntry := ProxyStruct{}

		err = json.Unmarshal(body, &proxyEntry)
		if err != nil {
			return fmt.Errorf("Error unmarshaling json response for credentials: %s", err)
		}

		for _, entry := range proxyEntry.Credentials.Entries {
			log.Println(entry.Name)
			log.Println(entry.Value)

			name := strings.ReplaceAll(entry.Name, " ", "")
			name = strings.Split(name, "(")[0]
			secrets[name] = []byte(entry.Value)
		}

	} else {
		return fmt.Errorf("code %d: %s", statusCode, body)
	}

	key := "/tmp/user.key"
	cert := "/tmp/user.cert"
	dest := "/tmp/" + user + ".pem"

	passwd := secrets["Passphrase"]

	if err := ioutil.WriteFile(key, secrets["PrivateKey"], os.FileMode(int(0400))); err != nil {
		return fmt.Errorf("Failed to write certificate key: %s", err)
	}

	if err := ioutil.WriteFile(cert, secrets["Certificate"], os.FileMode(int(0600))); err != nil {
		return fmt.Errorf("Failed to write certificate: %s", err)
	}

	err = CreateProxy(cert, key, string(passwd), dest)
	if err != nil {
		return fmt.Errorf("Failed to retrieve proxy: %s", err)
	}

	err = GetDN(dest)
	if err != nil {
		return fmt.Errorf("Failed to retrieve DN: %s", err)
	}

	return nil
}
