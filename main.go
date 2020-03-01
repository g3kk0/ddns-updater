package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"google.golang.org/api/dns/v1"
)

var errNotFound = errors.New("record not found")

func main() {
	//project := os.Getenv("PROJECT")
	//record := os.Getenv("RECORD")
	//zone := os.Getenv("ZONE")
	project := "crack-braid-160020"
	record := "foo.mkzd.host"
	zone := "mkzd-host"

	dnsService, err := newDnsService(project)
	if err != nil {
		log.Println(err)
	}

	dnsIp, err := dnsService.getRecordValue(zone, record)
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
		fmt.Println("updating DNS record")
		if dnsIp != "" {
			err := dnsService.deleteRecord(zone, record, dnsIp)
			if err != nil {
				log.Println(err)
			}
		}

		err = dnsService.updateRecord(zone, record, currentIp)
		if err != nil {
			log.Println(err)
		}
	} else {
		fmt.Println("record already up to date")
	}
}

type dnsService struct {
	ctx     context.Context
	client  *dns.Service
	project string
}

func newDnsService(project string) (*dnsService, error) {
	ctx := context.Background()
	client, err := dns.NewService(ctx)
	if err != nil {
		return nil, err
	}

	dnsService := &dnsService{
		ctx:     ctx,
		client:  client,
		project: project,
	}

	return dnsService, nil
}

func (s *dnsService) getRecordValue(zone, record string) (string, error) {
	resp, err := s.client.ResourceRecordSets.List(s.project, zone).Context(s.ctx).Do()
	if err != nil {
		return "", err
	}

	var ip string
	for _, v := range resp.Rrsets {
		if v.Name == record+"." {
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

	_, err := s.client.Changes.Create(s.project, zone, &change).Context(s.ctx).Do()
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

	_, err := s.client.Changes.Create(s.project, zone, &change).Context(s.ctx).Do()
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
