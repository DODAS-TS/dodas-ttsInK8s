package main

import (
	"flag"
	"log"

	tts "github.com/dodas-ts/dodas-ttsInK8s/pkg/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	configPath    string
	kubeClientset *kubernetes.Clientset
	cacheCerts    bool
	getProxy      bool
)

func init() {
	flag.BoolVar(&cacheCerts, "cache-certs", true, "Cache user credentials in k8s")
	flag.BoolVar(&getProxy, "get-proxy", false, "Get user proxy")
	flag.StringVar(&configPath, "config", ".config.yaml", "Path to yaml config file")

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error getting K8s in cluster config: %s", err)
	}

	kubeClientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating K8s clientset: %s", err)
	}

}

func main() {

	// TODO:
	// if --get-proxy take certSecret and to grid proxy init
	// if --get-certs take configMap and generate certSecret
	// do everything calling methods from package

	flag.Parse()

	var ttsClient *tts.TTSClient

	if cacheCerts {
		configFile := configPath

		ttsClient, err := ttsClient.GetConf(configFile)
		if err != nil {
			log.Fatalf("Error reading config file: %s", err)
		}

		err = ttsClient.CacheCredentials(kubeClientset * kubernetes.Clientset)
		if err != nil {
			log.Fatalf("Error caching credentials: %s", err)
		}
	} else if getProxy {

	}

}
