package main

import (
	"flag"
	"log"
	"time"

	tts "github.com/dodas-ts/dodas-ttsInK8s/pkg/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	configPath      string
	kubeClientset   *kubernetes.Clientset
	cacheCerts      bool
	getProxy        bool
	mapUser         bool
	dumpProxy       bool
	username        string
	token           string
	credentials     string
	periodInMinutes int
)

func init() {
	flag.BoolVar(&cacheCerts, "cache-certs", false, "Cache user credentials in k8s")
	flag.BoolVar(&getProxy, "get-proxy", false, "Get user proxy")
	flag.BoolVar(&mapUser, "map-user", false, "Get user proxy and map the DN to a username")
	flag.BoolVar(&dumpProxy, "dump-proxy", false, "Get proxy locally")
	flag.StringVar(&username, "user", "DUMMY", "Username for DN mapping")
	flag.StringVar(&token, "token", "", "Token for DN mapping")
	flag.StringVar(&credentials, "credentials", "https://dodas-tts.cloud.cnaf.infn.it/api/v2/iam/credential", "Specify IAM credentials endpoint")
	flag.StringVar(&configPath, "config", ".config.yaml", "Path to yaml config file")
	flag.IntVar(&periodInMinutes, "period", 120, "Proxy refresh period in minutes")

}

func main() {

	// TODO:
	// if --get-proxy take certSecret and to grid proxy init
	// if --get-certs take configMap and generate certSecret
	// do everything calling methods from package
	flag.Parse()

	if !dumpProxy {
		config, err := rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Error getting K8s in cluster config: %s", err)
		}

		kubeClientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("Error creating K8s clientset: %s", err)
		}
	}

	var ttsClient *tts.TTSClient

	if getProxy {
		for {
			err := tts.GetProxy("/root/proxy/gwms_proxy", kubeClientset)
			if err != nil {
				log.Printf("Error retrieving user proxy: %s", err)
			}
			time.Sleep(time.Duration(periodInMinutes * int(time.Minute)))
		}
	} else if cacheCerts {
		configFile := configPath
		ttsClient, err := ttsClient.GetConf(configFile)
		if err != nil {
			log.Fatalf("Error reading config file: %s", err)
		}

		err = ttsClient.CacheCredentials(kubeClientset)
		if err != nil {
			log.Fatalf("Error caching credentials: %s", err)
		}
	} else if mapUser {
		if token == "" || username == "" {
			log.Fatal("Please specify --username and --token")
		}

		configFile := configPath
		ttsClient, err := ttsClient.GetConf(configFile)
		if err != nil {
			log.Fatalf("Error reading config file: %s", err)
		}

		err = ttsClient.MapUser(username, token, kubeClientset)
		if err != nil {
			log.Fatalf("Failed to map %s with token %s: %s", username, token, err.Error())
		}
	} else if dumpProxy {
		if token == "" || username == "" {
			log.Fatal("Please specify --username and --token")
		}

		err := tts.DumpProxy(username, token, credentials)
		if err != nil {
			log.Fatalf("Error retrieving proxy: %s", err)
		}

	}

}
