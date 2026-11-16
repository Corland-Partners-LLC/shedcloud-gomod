package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Corland-Partners-LLC/shedcloud-gomod/partnerapi"
)

func main() {
	client, err := partnerapi.New(partnerapi.Options{
		// The ERP Partner API host is not built in: pass yours.
		BaseURL: os.Getenv("SHEDCLOUD_BASE_URL"),
		Auth:    partnerapi.Auth{APIKey: os.Getenv("SHEDCLOUD_API_KEY")},
	})
	if err != nil {
		log.Fatal(err)
	}

	page, err := client.StockUnits.List(context.Background(), partnerapi.StockUnitListParams{
		ListParams: partnerapi.ListParams{PageSize: 5},
		Statuses:   []string{"available"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d of %d stock units\n", len(page.Items), page.Pagination.Total)
	for _, u := range page.Items {
		fmt.Printf("%s  %-28s  %s  $%s\n", u.SerialNumber, u.ProductName, u.Status, u.AskingPrice)
	}
}
