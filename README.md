# mathemcli

CLI for searching products and managing your shopping cart on [Mathem.se](https://www.mathem.se) (Swedish grocery delivery).

## Installation

```bash
go install github.com/thepsadmin/mathemcli@latest
```

Or build from source:

```bash
git clone https://github.com/thepsadmin/mathemcli.git
cd mathemcli
go build -o mathemcli .
```

## Usage

### Login

```bash
mathemcli login
```

You'll be prompted for your Mathem email and password. Session is saved to `~/.mathemcli/session.json` and lasts ~30 days.

### Search Products

```bash
mathemcli search mjölk              # Search for milk
mathemcli search "arla ost"         # Multi-word search
mathemcli search kaffe --page 2     # See more results
```

### Manage Cart

```bash
mathemcli cart                  # View cart
mathemcli cart add 3681         # Add product by ID (from search)
mathemcli cart add 3681 3       # Add 3 of product
mathemcli cart clear            # Empty cart
```

### Logout

```bash
mathemcli logout
```

## Example Workflow

```bash
# Login
mathemcli login

# Search for milk
mathemcli search mellanmjölk
# [3681] ✓ Färsk Mellanmjölk 1,5%
#      Brand: Arla Ko®
#      1,5 l
#      Price: 19.95 SEK (13.30/l)

# Add to cart
mathemcli cart add 3681 2

# Check cart
mathemcli cart
# Cart: 2 varor (2 items)
# [3681] Arla Ko® Färsk Mellanmjölk 1,5%
#      Qty: 2 × 19.95 SEK = 39.90 SEK
```

## License

MIT
