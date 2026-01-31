package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/thepsadmin/mathemcli/internal/api"
)

var cartCmd = &cobra.Command{
	Use:   "cart",
	Short: "Manage shopping cart",
	Long:  `View and manage your Mathem shopping cart.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default action: show cart
		return showCart()
	},
}

var cartShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show cart contents",
	RunE: func(cmd *cobra.Command, args []string) error {
		return showCart()
	},
}

var cartAddCmd = &cobra.Command{
	Use:   "add [product_id] [quantity]",
	Short: "Add a product to cart",
	Long:  `Add a product to the cart by its ID. Get the ID from search results.`,
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		productID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid product ID: %w", err)
		}

		quantity := 1
		if len(args) > 1 {
			quantity, err = strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid quantity: %w", err)
			}
		}

		items := []api.CartItem{
			{ProductID: productID, Quantity: quantity},
		}

		cart, err := client.AddToCart(items)
		if err != nil {
			return fmt.Errorf("failed to add to cart: %w", err)
		}

		fmt.Printf("Added %d item(s) to cart\n", quantity)
		fmt.Printf("Cart total: %s %s (%d items)\n",
			cart.DisplayPrice, cart.Currency, cart.ProductQuantityCount)

		return nil
	},
}

var cartClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all items from cart",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := client.ClearCart()
		if err != nil {
			return fmt.Errorf("failed to clear cart: %w", err)
		}

		fmt.Println("Cart cleared")
		return nil
	},
}

func showCart() error {
	cart, err := client.GetCart()
	if err != nil {
		return fmt.Errorf("failed to get cart: %w", err)
	}

	if cart.ProductQuantityCount == 0 {
		fmt.Println("Your cart is empty")
		return nil
	}

	fmt.Printf("Cart: %s (%d items)\n\n", cart.LabelText, cart.ProductQuantityCount)

	for _, group := range cart.Groups {
		for _, item := range group.Items {
			fmt.Printf("[%d] %s\n", item.Product.ID, item.Product.FullName)
			if item.Product.NameExtra != "" {
				fmt.Printf("     %s\n", item.Product.NameExtra)
			}
			fmt.Printf("     Qty: %d × %s %s = %s %s\n",
				item.Quantity,
				item.Product.GrossPrice,
				cart.Currency,
				item.DisplayPrice,
				cart.Currency)
			fmt.Println()
		}
	}

	// Print summary
	fmt.Println("─────────────────────────────────")
	for _, summary := range cart.SummaryLines {
		for _, line := range summary.Lines {
			fmt.Printf("%-25s %s %s\n", line.Description, line.GrossAmount, cart.Currency)
		}
	}

	return nil
}

func init() {
	cartCmd.AddCommand(cartShowCmd)
	cartCmd.AddCommand(cartAddCmd)
	cartCmd.AddCommand(cartClearCmd)
}
