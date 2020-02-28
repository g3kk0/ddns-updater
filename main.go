package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"google.golang.org/api/dns/v1"
)

func main() {
	project := "crack-braid-160020"
	record := "foo.mkzd.host"
	zone := "mkzd-host"

	dnsService, err := newDnsService(project)
	if err != nil {
		log.Fatal(err)
	}

	dnsIp, err := dnsService.getRecordValue(record, zone)
	if err != nil {
		log.Fatal(err)
	}

	currentIp, err := getIpAddress()
	if err != nil {
		log.Fatal(err)
	}

	if dnsIp != currentIp {
		fmt.Println("updating DNS record")
	} else {
		fmt.Println("DNS record is in sync")
	}

	//fmt.Printf("ip = %+v\n", ip)

	// fetch ip
	// check

	//dnsService.updateRecord("foo", "bar")

	// fmt.Printf("ip = %+v\n", ip)

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

func (s *dnsService) getRecordValue(record, zone string) (string, error) {
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
		return "", errors.New("unable to find value in record")
	}

	return ip, nil
}

func (s *dnsService) updateRecord(name, ip string) error {
	change := dns.Change{
		Additions: []*dns.ResourceRecordSet{
			&dns.ResourceRecordSet{
				Kind: "dns#resourceRecordSet",
				Name: "foo.mkzd.host.",
				Rrdatas: []string{
					"1.2.3.4",
				},
				Ttl:  86400,
				Type: "A",
			},
		},
	}

	resp, err := s.client.Changes.Create("crack-braid-160020", "mkzd-host", &change).Context(s.ctx).Do()
	if err != nil {
		return err
	}

	spew.Dump(resp)

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
