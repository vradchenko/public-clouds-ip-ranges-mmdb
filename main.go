package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/maxmind/mmdbwriter"
	"github.com/maxmind/mmdbwriter/mmdbtype"
)

type Aws_ranges struct {
	SyncToken  string `json:"syncToken"`
	CreateDate string `json:"createDate"`
	Prefixes   []struct {
		IPPrefix           string `json:"ip_prefix"`
		Region             string `json:"region"`
		Service            string `json:"service"`
		NetworkBorderGroup string `json:"network_border_group"`
	} `json:"prefixes"`
	Ipv6Prefixes []struct {
		Ipv6Prefix         string `json:"ipv6_prefix"`
		Region             string `json:"region"`
		Service            string `json:"service"`
		NetworkBorderGroup string `json:"network_border_group"`
	} `json:"ipv6_prefixes"`
}

func main() {

	url := "https://ip-ranges.amazonaws.com/ip-ranges.json"

	writer, err := mmdbwriter.New(
		mmdbwriter.Options{
			DatabaseType: "AWS-IP-RANGES",
			RecordSize:   24,
		},
	)
	if err != nil {
		fmt.Println("Cannot start mmdbwriter")
		log.Fatal(err)
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Cannot download ip-ranges.json")
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Cannot read content of ip-ranges.json")
		log.Fatal(err)
	}

	var result Aws_ranges
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Cannot parse json from ip-ranges.json")
		log.Fatal(err)
	}

	for _, rec := range result.Prefixes {

		_, network, err := net.ParseCIDR(rec.IPPrefix)
		if err != nil {
			fmt.Println("ip_prefix is not valid")
			log.Fatal(err)
		}

		record := mmdbtype.Map{}
		record["region"] = mmdbtype.String(rec.Region)
		record["service"] = mmdbtype.String(rec.Service)
		record["network_border_group"] = mmdbtype.String(rec.NetworkBorderGroup)

		err = writer.Insert(network, record)
		if err != nil {
			fmt.Println("Cannot insert a record into mmdb file")
			log.Fatal(err)
		}

	}

	for _, rec := range result.Ipv6Prefixes {

		_, network, err := net.ParseCIDR(rec.Ipv6Prefix)
		if err != nil {
			log.Fatal(err)
			message := fmt.Sprintf("ParseCIDR error")
			return &message, err
		}

		record := mmdbtype.Map{}
		record["region"] = mmdbtype.String(rec.Region)
		record["service"] = mmdbtype.String(rec.Service)
		record["network_border_group"] = mmdbtype.String(rec.NetworkBorderGroup)

		err = writer.Insert(network, record)
		if err != nil {
			log.Fatal(err)
			message := fmt.Sprintf("Insert error")
			return &message, err
		}
	}	

	fh, err := os.Create("aws.mmdb")
	if err != nil {
		fmt.Println("Cannot create output file")
		log.Fatal(err)
	}

	_, err = writer.WriteTo(fh)
	if err != nil {
		fmt.Println("Cannot save output file")
		log.Fatal(err)
	}

}
