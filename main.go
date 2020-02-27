package main

import (
	"context"
	"log"

	"github.com/davecgh/go-spew/spew"
	"google.golang.org/api/dns/v1"
)

func main() {

	ctx := context.Background()
	dnsService, err := dns.NewService(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//resp, err := dnsService.ManagedZones.Get("crack-braid-160020", "mkzd-host").Context(ctx).Do()

	// Get record
	//resp, err := dnsService.ResourceRecordSets.List("crack-braid-160020", "mkzd-host").Context(ctx).Do()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Update record
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

	resp, err := dnsService.Changes.Create("crack-braid-160020", "mkzd-host", &change).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(resp)
}
