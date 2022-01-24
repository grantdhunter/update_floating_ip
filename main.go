package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/digitalocean/godo"
)

func get_ip(ctx context.Context, client *godo.Client, domain_name string) (*godo.FloatingIP, error) {
	record, _, err := client.Domains.RecordsByType(ctx, domain_name, "A", nil)
	if err != nil {
		return nil, err
	}

	floating_ip, _, err := client.FloatingIPs.Get(ctx, record[0].Data)

	if err != nil {
		return nil, err
	}

	return floating_ip, nil

}

func get_droplets_by_tag(ctx context.Context, client *godo.Client, tag string) ([]godo.Droplet, error) {
	droplets, _, err := client.Droplets.ListByTag(ctx, tag, &godo.ListOptions{})
	if err != nil {
		return nil, err
	}

	return droplets, nil
}

func assign_ip_to_droplet(ctx context.Context, client *godo.Client, ip string, id int) (*godo.Action, error) {
	action, _, err := client.FloatingIPActions.Assign(ctx, ip, id)
	if err != nil {
		return nil, err
	}

	return action, nil
}

var (
	wait bool
)

func init() {
	flag.BoolVar(&wait, "wait", false, "Flag to have the program wait indefinitely.")
}

func main() {

	token := os.Getenv("DO_API_TOKEN")
	if len(token) == 0 {
		log.Println("DO_API_TOKEN not set.")
		os.Exit(1)
	}

	client := godo.NewFromToken(token)
	ctx := context.TODO()

	flag.Parse()
	domain := flag.Arg(0)

	if len(domain) == 0 {
		log.Println("Usage: main.go [-wait] domain")
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Printf("Querying domain: %s\n", domain)
	ip, err := get_ip(ctx, client, domain)
	if err != nil {
		log.Printf("Failure: %s\n", err)
		os.Exit(1)
	}

	log.Printf("Current IP: %s\n", *&ip.IP)

	if ip.Droplet != nil {
		log.Printf("Current droplet: %s (%d)\n", ip.Droplet.Name, ip.Droplet.ID)
	} else {

		droplets, err := get_droplets_by_tag(ctx, client, "k8s:worker")
		if err != nil {
			log.Printf("Failure: %s\n", err)
			os.Exit(1)

		}
		if len(droplets) == 0 {
			log.Println("No Droplets found")
			os.Exit(1)
		}

		_, err = assign_ip_to_droplet(ctx, client, ip.IP, droplets[0].ID)

		if err != nil {
			log.Printf("Failure: %s\n", err)
			os.Exit(1)

		}
		log.Printf("Assigned droplet %s to IP %s\n", droplets[0].Name, ip.IP)
	}

	if wait {
		log.Println("Wait activated, sleeping until killed")

		for {
			log.Println("beep")
			time.Sleep(time.Hour * 24)
		}
	}

}
