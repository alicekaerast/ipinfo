package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path"
)

type AccountData struct {
	Accounts []struct {
		Cidr        string `yaml:"cidr"`
		Account     string `yaml:"account"`
		Description string `yaml:"description,omitempty"`
	} `yaml:"accounts"`
}

func main() {
	app := &cli.App{
		Action: func(c *cli.Context) error {
			wantedIPString := c.Args().Get(0)
			wantedIP := net.ParseIP(wantedIPString)
			fmt.Printf("Searching %v\n", wantedIP)

			var accountData AccountData
			home, _ := os.UserHomeDir()
			yamlFile, err := ioutil.ReadFile(path.Join(home, ".accounts.yml"))
			if err != nil {
				return err
			}
			err = yaml.Unmarshal(yamlFile, &accountData)
			if err != nil {
				return err
			}
			for _, account := range accountData.Accounts {
				_, network, err := net.ParseCIDR(account.Cidr)
				if err != nil {
					return nil
				}

				if network.Contains(wantedIP) {
					fmt.Printf("Found Account: %v %q %v\n", account.Account, account.Description, account.Cidr)
					os.Setenv("AWS_PROFILE", account.Account)
				}
			}

			sess := session.Must(session.NewSession(&aws.Config{Region: aws.String("eu-west-1")}))
			ec2Client := ec2.New(sess)

			filters := make([]*ec2.Filter, 0)
			keyName := "addresses.private-ip-address"
			filter := ec2.Filter{
				Name: &keyName, Values: []*string{&wantedIPString}}
			filters = append(filters, &filter)
			dnii := &ec2.DescribeNetworkInterfacesInput{Filters: filters}

			resp, err := ec2Client.DescribeNetworkInterfaces(dnii)
			networkInterfaces := resp.NetworkInterfaces
			if len(networkInterfaces) == 0 {
				fmt.Println("Network interface not found")
				os.Exit(1)
			}
			description := *resp.NetworkInterfaces[0].Description
			if description != "" {
				fmt.Printf("Description: %q\n", description)
			}
			attachment := *resp.NetworkInterfaces[0].Attachment
			if attachment.InstanceId != nil {
				fmt.Printf("Instance: %v\n", attachment.InstanceId)
				input := &ec2.DescribeInstancesInput{
					InstanceIds: []*string{
						aws.String(*attachment.InstanceId),
					},
				}
				instances, _ := ec2Client.DescribeInstances(input)
				tags := instances.Reservations[0].Instances[0].Tags
				for _, j := range tags {
					if *j.Key == "Name" {
						fmt.Printf("Name: %v\n", *j.Value)
					}
				}
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
