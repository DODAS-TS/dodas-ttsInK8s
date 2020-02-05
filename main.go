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
	periodInMinutes int
)

func init() {
	flag.BoolVar(&cacheCerts, "cache-certs", false, "Cache user credentials in k8s")
	flag.BoolVar(&getProxy, "get-proxy", false, "Get user proxy")
	flag.StringVar(&configPath, "config", ".config.yaml", "Path to yaml config file")
	flag.IntVar(&periodInMinutes, "period", 120, "Proxy refresh period in minutes")

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

	if getProxy {
		for {
			err := tts.GetProxy("/root/gwms_proxy", kubeClientset)
			if err != nil {
				log.Printf("Error retrieving user proxy: %s", err)
			}
			time.Sleep(periodInMinutes * time.Minute)
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
	}

}
