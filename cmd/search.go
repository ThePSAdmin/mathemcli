package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	searchPage int
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for products",
	Long:  `Search for products on Mathem by name or keyword.`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.Join(args, " ")

		result, err := client.Search(query, searchPage)
		if err != nil {
			return fmt.Errorf("search failed: %w", err)
		}

		if len(result.Items) == 0 {
			fmt.Println("No products found")
			return nil
		}

		fmt.Printf("Found %d products (page %d):\n\n",
			result.Attributes.Items, result.Attributes.Page)

		for _, item := range result.Items {
			if item.Type != "product" {
				continue
			}
			attr := item.Attributes
			availability := "✓"
			if !attr.Availability.IsAvailable {
				availability = "✗"
			}

			fmt.Printf("[%d] %s %s\n", item.ID, availability, attr.Name)
			if attr.Brand != "" {
				fmt.Printf("     Brand: %s\n", attr.Brand)
			}
			if attr.NameExtra != "" {
				fmt.Printf("     %s\n", attr.NameExtra)
			}
			fmt.Printf("     Price: %s %s", attr.GrossPrice, attr.Currency)
			if attr.GrossUnitPrice != "" && attr.UnitPriceQuantityAbbr != "" {
				fmt.Printf(" (%s/%s)", attr.GrossUnitPrice, attr.UnitPriceQuantityAbbr)
			}
			fmt.Println()
			fmt.Println()
		}

		if result.Attributes.HasMoreItems {
			fmt.Printf("More results available. Use --page %d to see next page.\n",
				searchPage+1)
		}

		return nil
	},
}

func init() {
	searchCmd.Flags().IntVarP(&searchPage, "page", "n", 1, "Page number")
}
