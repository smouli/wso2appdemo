package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var (
	oauthConf = &clientcredentials.Config{
		ClientID:/*"fy2TvqkILON1nfsqL6zaLL6C0m4a"*/ "HCZAOMKxHCeHfuYwmj7iuRPMbpEa",
		ClientSecret:/*"9Aon6fkYoEBeGBawwg4fWbqpg6Aa"*/ "WxgaU5jG2dhmjbIQKuy6zHZw44Ya",
		TokenURL: "https://localhost:9443/oauth2/token",
	}
	// random string for oauth2 API calls to protect against CSRF
	oauthStateString = "thisshouldberandom"
)

func getAccessToken() (string, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	token, err := oauthConf.Token(oauth2.NoContext)
	if err != nil {
		fmt.Println("Error In Token Request")
		return "", err
	}
	return token.AccessToken, nil
}

// func main() {
// 	accessToken, err := getAccessToken()
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	fmt.Println(accessToken)
// }

// func main() {
// 	IDPTokenUrl := "https://localhost:9443"
// 	resource := "/oauth2/token"
// 	u, err := url.ParseRequestURI(IDPTokenUrl)
// 	if err != nil {
// 		//return nil, err
// 	}

// 	u.Path = resource
// 	urlStr := u.String()
// 	data := url.Values{}
// 	data.Set("grant_type", "client_credentials")
// 	data.Add("validity_period", "3600")
// 	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
// 	client := &http.Client{}
// 	r, err := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode()))
// 	if err != nil {
// 		//return nil, err
// 	}

// 	clientCredArr := []string{oauthConf.ClientID, oauthConf.ClientSecret}
// 	clientCredStr := strings.Join(clientCredArr, ":")
// 	encrCred := base64.StdEncoding.EncodeToString([]byte(clientCredStr))
// 	encrCred = strings.Join([]string{"Basic", encrCred}, " ")

// 	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 	r.Header.Add("Authorization", encrCred)
// 	resp, err := client.Do(r)
// 	if err != nil {
// 		fmt.Println("ERROR")
// 		//return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var t Token
// 	err = json.Unmarshal(resp.Body, &t)
// 	if err != nil {
// 		fmt.Println("json.Unmarshal", err)
// 	}
// 	fmt.Printf("Token:\n")
// 	fmt.Println(t)

// }
