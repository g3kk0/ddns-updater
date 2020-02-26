package main

import (
	"context"
	"fmt"
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

	resp, err := dnsService.ManagedZones.Get(project, managedZone).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("dnsService.ManagedZones = %+v\n", dnsService.ManagedZones)

	spew.Dump(resp)

	// ctx := context.Background()

	// c, err := google.DefaultClient(ctx, dns.CloudPlatformScope)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// dnsService, err := dns.New(c)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Identifies the project addressed by this request.
	// project := "my-project" // TODO: Update placeholder value.

	// // Identifies the managed zone addressed by this request. Can be the managed zone name or id.
	// managedZone := "my-managed-zone" // TODO: Update placeholder value.

	// resp, err := dnsService.ManagedZones.Get(project, managedZone).Context(ctx).Do()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // TODO: Change code below to process the `resp` object:
	// fmt.Printf("%#v\n", resp)
	fmt.Println("done")
}
