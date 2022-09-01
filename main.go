package main

import (
	"fmt"
	go_rapi "github.com/sarff/go-rapi"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -99
}

type NetPlan struct {
	Network struct {
		Version   int `yaml:"version"`
		Ethernets struct {
			Eth0 struct {
				Addresses []string `yaml:"addresses"`
				Gateway4  string   `yaml:"gateway4"`
				Gateway6  string   `yaml:"gateway6"`
				Match     struct {
					Macaddress string `yaml:"macaddress"`
				} `yaml:"match"`
				Nameservers struct {
					Addresses []string `yaml:"addresses"`
					Search    []string `yaml:"search"`
				} `yaml:"nameservers"`
				SetName string `yaml:"set-name"`
			} `yaml:"eth0"`
		} `yaml:"ethernets"`
	} `yaml:"network"`
}

func main() {
	iplist := []string{"91.20.19.22/24", "91.208.19.21/24", "91.20.19.5/24"}
	path := "/etc/netplan/50-cloud-init.yaml"
	zone := "686cc90e5031d8*****87ae4"
	auth_key := "7981626eb9c5*****8c69a2c8e3a9dce"
	dns_id := "20482e9887a179*****61792241e6"
	auth_mail := "cloudflare@example.com"
	domein_name := "mail.example.com"
	//path := "./50-cloud-init.yaml"

	file, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	np := NetPlan{}

	err = yaml.Unmarshal(file, &np)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	indx := Find(iplist, np.Network.Ethernets.Eth0.Addresses[0])
	curlput := ""
	if indx < 2 {
		np.Network.Ethernets.Eth0.Addresses[0] = iplist[indx+1]
		curlput = strings.Split(iplist[indx+1], "/")[0]
	} else {
		np.Network.Ethernets.Eth0.Addresses[0] = iplist[0]
		curlput = strings.Split(iplist[0], "/")[0]
	}

	res, err := yaml.Marshal(&np)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(path, res, 0777)
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command("netplan", "apply")

	params := make(map[string]interface{})
	params["type"] = "A"
	params["name"] = domein_name
	params["content"] = curlput
	params["ttl"] = "1"
	params["proxied"] = true

	headers := map[string]string{"X-Auth-Email": auth_mail, "X-Auth-Key": auth_key, "Content-Type": "application/json"}

	_, err = go_rapi.HttpPutBody(fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", zone, dns_id), params,
		headers)

	if err != nil {
		log.Fatal(err)
	}

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

}
