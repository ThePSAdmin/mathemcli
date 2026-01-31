---
name: mathemcli
description: Use when user wants to search Mathem grocery products, manage shopping cart, or interact with Mathem.se via CLI
---

# mathemcli

CLI for searching products and managing shopping cart on Mathem.se (Swedish grocery delivery).

## Quick Reference

| Command | Description |
|---------|-------------|
| `mathemcli login` | Authenticate (prompts for email/password) |
| `mathemcli logout` | Clear saved session |
| `mathemcli search <query>` | Search products by name |
| `mathemcli cart` | Show cart contents |
| `mathemcli cart add <id> [qty]` | Add product to cart |
| `mathemcli cart clear` | Empty the cart |

## Authentication

Login is required before using search or cart commands:

```bash
mathemcli login
# Email: user@example.com
# Password: (hidden)
```

Session is saved to `~/.mathemcli/session.json` and lasts ~30 days.

## Search Products

```bash
mathemcli search mjölk           # Search for milk
mathemcli search "arla ost"      # Multi-word search
mathemcli search kaffe --page 2  # Pagination
```

Output shows product ID (needed for cart), availability, name, brand, size, and price.

## Cart Management

```bash
# View cart
mathemcli cart

# Add items (use product ID from search)
mathemcli cart add 3681      # Add 1 item
mathemcli cart add 3681 3    # Add 3 items

# Clear cart
mathemcli cart clear
```

## Typical Workflow

```bash
mathemcli login
mathemcli search "mellanmjölk"   # Find product, note ID
mathemcli cart add 3681 2        # Add 2 to cart
mathemcli cart                   # Verify cart
```

## Common Issues

| Issue | Solution |
|-------|----------|
| "not logged in" | Run `mathemcli login` |
| 403 errors | Session expired, login again |
| Product not found | Use Swedish search terms |
