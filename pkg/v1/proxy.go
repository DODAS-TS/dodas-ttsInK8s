package v1

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	execute "github.com/alexellis/go-execute/pkg/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetProxy ..
func GetProxy(dest string, kubeClientset *kubernetes.Clientset) error {

	secretCert, err := kubeClientset.CoreV1().Secrets(v1.NamespaceDefault).Get("certs-secret", metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Problem checking for cert secret existance: %s", err)
	}

	key := "/tmp/user.key"
	cert := "/tmp/user.cert"

	passwd := secretCert.Data["Passphrase"]

	if err := ioutil.WriteFile(key, secretCert.Data["PrivateKey"], os.FileMode(int(0400))); err != nil {
		return fmt.Errorf("Failed to write certificate key: %s", err)
	}

	if err := ioutil.WriteFile(cert, secretCert.Data["Certificate"], os.FileMode(int(0600))); err != nil {
		return fmt.Errorf("Failed to write certificate: %s", err)
	}

	log.Print("Generating user proxy")

	cmd := execute.ExecTask{
		Command: "echo",
		Args: []string{
			string(passwd),
			"|",
			"grid-proxy-init",
			"-valid",
			"160:00",
			"-key",
			key,
			"-cert",
			cert,
			"-out",
			dest,
		},
		StreamStdio: false,
		Shell:       true,
	}

	res, err := cmd.Execute()
	if err != nil {
		return fmt.Errorf("grid-proxy-init command failed: %s", err)
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("grid-proxy-init exited with error: %s", res.Stderr)
	}

	log.Printf("stdout: %s, stderr: %s, exit-code: %d\n", res.Stdout, res.Stderr, res.ExitCode)

	return nil
}
