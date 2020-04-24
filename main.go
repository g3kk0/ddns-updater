package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"google.golang.org/api/dns/v1"
)

var errNotFound = errors.New("DNS record not found")

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	config := loadConfig()

	dnsService, err := newDnsService(config.googleProjectId)
	if err != nil {
		return err
	}

	log.Println("Starting ddns-updater...")
	log.Printf("Check interval = %s\n", config.checkInterval)

	for {
		dnsIp, err := dnsService.getRecordValue(config.dnsZone, config.dnsRecord)
		if err != nil {
			switch err {
			case errNotFound:
				log.Println(errNotFound)
			default:
				log.Println(err)
			}
		}

		currentIp, err := getIpAddress()
		if err != nil {
			log.Println(err)
		}

		if dnsIp != currentIp {
			log.Printf("Updating DNS record %s with IP %s\n", config.dnsRecord, currentIp)
			if dnsIp != "" {
				err := dnsService.deleteRecord(config.dnsZone, config.dnsRecord, dnsIp)
				if err != nil {
					log.Println(err)
				}
			}

			err = dnsService.updateRecord(config.dnsZone, config.dnsRecord, currentIp)
			if err != nil {
				log.Println(err)
			} else {
				log.Println("DNS record sucessfully updated")
			}
		} else {
			log.Printf("DNS record up to date (%s -> %s)\n", config.dnsRecord, currentIp)
		}

		interval, err := strconv.Atoi(config.checkInterval)
		if err != nil {
			log.Println(err)
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}

type config struct {
	googleProjectId string
	dnsRecord       string
	dnsZone         string
	checkInterval   string
}

func loadConfig() *config {
	config := &config{}

	config.googleProjectId = os.Getenv("GOOGLE_PROJECT_ID")
	config.dnsRecord = os.Getenv("DNS_RECORD")
	config.dnsZone = os.Getenv("DNS_ZONE")
	config.checkInterval = os.Getenv("CHECK_INTERVAL")

	if config.checkInterval == "" {
		config.checkInterval = "3600"
	}

	return config
}

type dnsService struct {
	ctx       context.Context
	client    *dns.Service
	projectId string
}

func newDnsService(projectId string) (*dnsService, error) {
	ctx := context.Background()
	client, err := dns.NewService(ctx)
	if err != nil {
		return nil, err
	}

	dnsService := &dnsService{
		ctx:       ctx,
		client:    client,
		projectId: projectId,
	}

	return dnsService, nil
}

func (s *dnsService) getRecordValue(zone, record string) (string, error) {
	resp, err := s.client.ResourceRecordSets.List(s.projectId, zone).Context(s.ctx).Do()
	if err != nil {
		return "", err
	}

	var ip string
	for _, v := range resp.Rrsets {
		if v.Name == record+"." && v.Type == "A" {
			ip = v.Rrdatas[0]
		}
	}

	if ip == "" {
		return "", errNotFound
	}

	return ip, nil
}

func (s *dnsService) deleteRecord(zone, record, ip string) error {
	change := dns.Change{
		Deletions: []*dns.ResourceRecordSet{
			&dns.ResourceRecordSet{
				Kind: "dns#resourceRecordSet",
				Name: record + ".",
				Rrdatas: []string{
					ip,
				},
				Ttl:  300,
				Type: "A",
			},
		},
	}

	_, err := s.client.Changes.Create(s.projectId, zone, &change).Context(s.ctx).Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *dnsService) updateRecord(zone, record, ip string) error {
	change := dns.Change{
		Additions: []*dns.ResourceRecordSet{
			&dns.ResourceRecordSet{
				Kind: "dns#resourceRecordSet",
				Name: record + ".",
				Rrdatas: []string{
					ip,
				},
				Ttl:  300,
				Type: "A",
			},
		},
	}

	_, err := s.client.Changes.Create(s.projectId, zone, &change).Context(s.ctx).Do()
	if err != nil {
		return err
	}

	return nil
}

func getIpAddress() (string, error) {
	resp, err := http.Get("https://myexternalip.com/raw")
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return string(body), nil
}
