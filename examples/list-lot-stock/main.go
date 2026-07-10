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
		// Environment: partnerapi.EnvironmentSandbox, // optional
		Auth: partnerapi.Auth{
			APIKey: os.Getenv("SHEDCLOUD_API_KEY"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	stock, err := client.LotStock.List(context.Background(), partnerapi.LotStockListParams{
		PaginationParams: partnerapi.PaginationParams{Limit: 5},
		PurchaseType:     "Lot Stock",
		Sort:             "price",
		Order:            partnerapi.SortAsc,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("total=%d returned=%d\n", stock.Total, len(stock.Data))
	if len(stock.Data) > 0 {
		item := stock.Data[0]
		fmt.Printf("first: %s | %s | $%.2f\n", item.SerialNumber, item.Title, item.Price)
	}
}
