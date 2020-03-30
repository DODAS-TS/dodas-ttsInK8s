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

// GetDN ..
func GetDN(proxy string) error {
	cmd := execute.ExecTask{
		Command: "voms-proxy-info",
		Args: []string{
			"--file",
			proxy,
			"--issuer",
		},
		StreamStdio: false,
		Shell:       true,
	}

	res, err := cmd.Execute()
	if err != nil {
		return fmt.Errorf("voms-proxy-info command failed: %s", err)
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("voms-proxy-info exited with error: %s", res.Stderr)
	}

	log.Printf("UserDN: %s\n", res.Stdout)

	return nil
}

// CreateProxy ..
func CreateProxy(cert string, key string, passwd string, dest string) error {

	cmd := execute.ExecTask{
		Command: "echo",
		Args: []string{
			string(passwd),
			"|",
			"voms-proxy-init",
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
		return fmt.Errorf("voms-proxy-init command failed: %s", err)
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("voms-proxy-init exited with error: %s", res.Stderr)
	}

	log.Printf("stdout: %s, stderr: %s, exit-code: %d\n", res.Stdout, res.Stderr, res.ExitCode)

	log.Printf("checking certs validity.")

	//voms-proxy-info --valid 5:00 -e --file /root/proxy/gwms_proxy
	//openssl x509 -checkend 86400 -noout -in file.pem
	cmd = execute.ExecTask{
		Command: "openssl",
		Args: []string{
			"x509",
			"-checkend",
			"18000",
			"-noout",
			"-in",
			cert,
		},
		StreamStdio: false,
		Shell:       true,
	}

	res, err = cmd.Execute()
	if err != nil {
		return fmt.Errorf("openssl command failed: %s", err)
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("openssl exited with error: %s", res.Stderr)
	}

	log.Printf("stdout: %s, stderr: %s, exit-code: %d\n", res.Stdout, res.Stderr, res.ExitCode)

	return nil
}

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

	err = CreateProxy(cert, key, string(passwd), dest)
	if err != nil {
		return fmt.Errorf("Failed to retrieve proxy: %s", err)
	}

	return nil
}
