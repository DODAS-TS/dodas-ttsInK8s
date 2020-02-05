package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

// GetAccessToken ..
func (clientConf Conf) GetAccessToken(refreshToken string) (token string, err error) {

	clientID := clientConf.AllowRefresh.ClientID
	clientSecret := clientConf.AllowRefresh.ClientSecret
	IAMTokenEndpoint := clientConf.AllowRefresh.IAMTokenEndpoint

	v := url.Values{}

	v.Set("client_id", clientID)
	v.Set("client_secret", clientSecret)
	v.Set("grant_type", "refresh_token")
	v.Set("refresh_token", refreshToken)

	request := Request{
		URL:         IAMTokenEndpoint,
		RequestType: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		AuthUser: clientID,
		AuthPwd:  clientSecret,
		Content:  []byte(v.Encode()),
	}

	body, statusCode, err := MakeRequest(request)
	if err != nil {
		return "", err
	}

	if statusCode != 200 {
		fmt.Println("ERROR:\n", string(body))
		return "", fmt.Errorf("Error code %d: %s", statusCode, string(body))
	}

	var bodyJSON RefreshTokenStruct

	//fmt.Println(string(body))
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		return "", err
	}

	// TODO: only if the mode is IAM for both cloud and
	clientConf.Cloud.Password = bodyJSON.AccessToken
	clientConf.Im.Token = bodyJSON.AccessToken

	return bodyJSON.AccessToken, nil
}

// GetNewToken ..
func (clientConf Conf) GetNewToken() (updatedConf Conf, err error) {

	tokenBytes, err := ioutil.ReadFile(clientConf.AllowRefresh.RefreshTokenFile)
	if err != nil {
		return Conf{}, fmt.Errorf("Failed to read refresh token, please be sure you did `dodas iam init` command: %s", err)
	}

	accessToken, err := clientConf.GetAccessToken(string(tokenBytes))
	if err != nil {
		return Conf{}, err
	}

	//fmt.Printf("Access token: %s", accessToken)

	// TODO: only if the mode is IAM for both cloud and
	clientConf.Cloud.Password = accessToken
	clientConf.Im.Token = accessToken

	// TODO: dump access token to a file

	updatedConf = clientConf
	return updatedConf, nil
}

// GetRefreshToken ..
func (clientConf Conf) GetRefreshToken() (RefreshToken string, err error) {

	clientID := clientConf.AllowRefresh.ClientID
	clientSecret := clientConf.AllowRefresh.ClientSecret
	IAMTokenEndpoint := clientConf.AllowRefresh.IAMTokenEndpoint
	accessToken := clientConf.Im.Token

	v := url.Values{}

	v.Set("client_id", clientID)
	v.Set("client_secret", clientSecret)
	v.Set("grant_type", "urn:ietf:params:oauth:grant-type:token-exchange")
	v.Set("subject_token", accessToken)

	request := Request{
		URL:         IAMTokenEndpoint,
		RequestType: "POST",
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
		AuthUser: clientID,
		AuthPwd:  clientSecret,
		Content:  []byte(v.Encode()),
	}

	body, statusCode, err := MakeRequest(request)
	if err != nil {
		return "", err
	}

	if statusCode != 200 {
		fmt.Printf("Error code %d: %s\n", statusCode, string(body))
		return "", fmt.Errorf("Error code %d: %s", statusCode, string(body))
	}

	var bodyJSON RefreshTokenStruct

	//fmt.Println(string(body))
	err = json.Unmarshal(body, &bodyJSON)
	if err != nil {
		return "", err
	}

	return bodyJSON.RefreshToken, nil
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a IAM context for automatic token refresh",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		isTokenUsed := (clientConf.Im.Token != "" || clientConf.Cloud.AuthVersion == "3.x_oidc_access_token")
		isRefreshSet := clientConf.AllowRefresh.IAMTokenEndpoint != ""

		if !isTokenUsed && !isRefreshSet {
			panic(fmt.Errorf("Token not used anywhere in config or refresh endpoint missing"))
		}

		fmt.Println("Removing old dump files and creating new ones")

		// Remove dumps files if exists
		if _, err := os.Stat(clientConf.AllowRefresh.AccessTokenFile); err == nil {
			err := os.Remove(clientConf.AllowRefresh.AccessTokenFile)
			if err != nil {
				panic(err)
			}
		}

		if _, err := os.Stat(clientConf.AllowRefresh.RefreshTokenFile); err == nil {
			err = os.Remove(clientConf.AllowRefresh.RefreshTokenFile)
			if err != nil {
				panic(err)
			}
		}

		if clientConf.Im.Token == "" {
			panic(fmt.Errorf("Error: access token not specified to IM"))
		}

		if err := ioutil.WriteFile(clientConf.AllowRefresh.AccessTokenFile, []byte(clientConf.Im.Token), os.FileMode(int(0600))); err != nil {
			log.Fatal(err)
		}

		token, err := clientConf.GetRefreshToken()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Got refresh token: %s", token)

		if err := ioutil.WriteFile(clientConf.AllowRefresh.RefreshTokenFile, []byte(token), os.FileMode(int(0600))); err != nil {
			log.Fatal(err)
		}

	},
}

// iamCmd represents the iam command
var iamCmd = &cobra.Command{
	Use:   "iam",
	Short: "Wrapper command for IAM interaction",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("iam called")
	},
}

func init() {
	rootCmd.AddCommand(iamCmd)
	iamCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// iamCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// iamCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
