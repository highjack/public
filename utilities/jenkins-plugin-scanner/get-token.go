package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/base64"
	"encoding/json"
	"strings"
	"bufio"
	"os"
)
type CrumbData struct {
	Class              string `json:"_class"`
	Crumb              string `json:"crumb"`
	CrumbRequestField string `json:"crumbRequestField"`
}

type TokenData struct {
	TokenName  string `json:"tokenName"`
	TokenUUID  string `json:"tokenUuid"`
	TokenValue string `json:"tokenValue"`
}

type Response struct {
	Status string    `json:"status"`
	Data   TokenData `json:"data"`
}

func main() {
	url := "http://localhost:8080/crumbIssuer/api/json"

	username := "admin"

	file, err := os.Open("/var/jenkins_home/secrets/initialAdminPassword")
	//for testing locally: file, err := os.Open("/tmp/pass")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	password := ""
	if scanner.Scan() {
		password = scanner.Text()
	} else if err := scanner.Err(); err != nil {
	fmt.Println("unable to read password file")
		os.Exit(0)

	}
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", authHeader)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var crumbData CrumbData
	err = json.Unmarshal([]byte(body), &crumbData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	cookieValue := ""
	for _, cookie := range resp.Cookies() {
		cookieValue = fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
	}


	crumbValue := crumbData.Crumb

	postURL := "http://localhost:8080/user/admin/descriptorByName/jenkins.security.ApiTokenProperty/generateNewToken"

	req, err = http.NewRequest("POST", postURL, strings.NewReader("newTokenName=foo"))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Cookie", cookieValue)
	req.Header.Set("Jenkins-Crumb", crumbValue)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var response Response

	// Unmarshal the JSON data into the Response struct
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	config := ("current: yourServer\n" +
		  "jenkins_servers:\n" +
		  "- name: yourServer\n" +
  		  " url: http://localhost:8080\n" +
		  " username: admin\n" +
		  " token: " + response.Data.TokenValue + "\n" +
  		  " insecureSkipVerify: true\n" +
		  "mirrors:\n" + 
		  "- name: default\n" + 
 		  " url: http://mirrors.jenkins.io/\n" +
		  "- name: tsinghua\n" +
		  " url: https://mirrors.tuna.tsinghua.edu.cn/jenkins/\n"+
		  "- name: huawei\n"+
		  " url: https://mirrors.huaweicloud.com/jenkins/\n"+
		  "- name: tencent"+
		  " url: https://mirrors.cloud.tencent.com/jenkins/\n")
	

	  dirname, err := os.UserHomeDir()
    	if err != nil {
        	fmt.Println( err )
    	}

	f, err := os.Create(dirname+"/.jenkins-cli.yaml")
	if err != nil{
		fmt.Println(err)
		return
	}
	_, err = f.WriteString(config)
	if err != nil{
		fmt.Println("unable to write config")
		return
		}
	f.Sync()
	fmt.Println(config)

}

