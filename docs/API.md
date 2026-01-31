# Mathem API Documentation

This document describes the undocumented Mathem REST API discovered through browser inspection.

## Base URL

```
https://www.mathem.se/tienda-web-api/v1/
```

## Authentication

Mathem uses session-based authentication via cookies.

### Login

**Endpoint:** `POST /user/login/`

**Request:**
```json
{
  "username": "user@example.com",
  "password": "your-password"
}
```

**Response:** Sets `sessionid` cookie (httpOnly, secure, SameSite=Lax)

**Session Duration:** ~30 days

### Making Authenticated Requests

Include the session cookie in all requests:
```
Cookie: sessionid=<session_value>
Cookie: csrftoken=<csrf_value>
```

### Bot Protection

The API has bot protection that requires:

1. **Browser-like User-Agent**: Use a real Chrome/Firefox user-agent string
2. **Client Hints Headers**: Include `sec-ch-ua`, `sec-ch-ua-mobile`, `sec-ch-ua-platform`
3. **Referer Header**: Include a valid referer from `https://www.mathem.se`
4. **Initial Page Visit**: Before login, visit `/se/user/login/` to get the initial `csrftoken` cookie

Example headers:
```
User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36
sec-ch-ua: "Not(A:Brand";v="8", "Chromium";v="144", "Google Chrome";v="144"
sec-ch-ua-mobile: ?0
sec-ch-ua-platform: "Linux"
Referer: https://www.mathem.se/se/
```

## Endpoints

### Search

#### Search Products

**Endpoint:** `GET /search/mixed/`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `q` | string | Search query |
| `type` | string | `product` for products, `suggestion` for autocomplete |
| `page` | int | Page number (default: 1) |
| `items` | int | Items per page |

**Example:**
```
GET /search/mixed/?q=mjölk&type=product&page=1
```

**Response:** Returns product list with details including:
- `id` - Product ID (used for cart operations)
- `attributes.name` - Product name
- `attributes.price` - Price information

### Cart

#### Get Cart

**Endpoint:** `GET /cart/`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| `group_by` | string | `recipes` or `categories` |

**Example:**
```
GET /cart/?group_by=recipes
```

**Response:**
```json
{
  "id": 0,
  "product_quantity_count": 2,
  "display_price": "278.00",
  "total_gross_amount": "411.80",
  "currency": "SEK",
  "groups": [
    {
      "items": [
        {
          "product": {
            "id": 2352,
            "full_name": "Arla Gouda 28%",
            "gross_price": "139.00"
          },
          "item_id": 97907308,
          "quantity": 2
        }
      ]
    }
  ]
}
```

#### Add Items to Cart

**Endpoint:** `POST /cart/items/`

**Request:**
```json
{
  "items": [
    {
      "product_id": 2352,
      "quantity": 1
    }
  ]
}
```

**Response:** Returns updated cart state

**Note:** Quantities are additive. To set a specific quantity, calculate the delta.

#### Clear Cart

**Endpoint:** `POST /cart/clear/`

**Request:** Empty body

**Response:** Returns empty cart state

### Other Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/dixa/user-jwt/` | GET | Get user JWT for support chat |
| `/slot-picker/slots/?num-days=3` | GET | Get delivery time slots |
| `/cart/validate/` | GET | Validate cart contents |
| `/campaigns/promoted_products/` | GET | Get promoted products |
| `/perks/` | GET | Get user perks/rewards |
| `/app-components/home/` | GET | Get homepage components |

## Error Handling

Errors return JSON with the following structure:
```json
{
  "errors": ["Error message"],
  "field_errors": {
    "field_name": {
      "message": "Field required",
      "code": "missing"
    }
  }
}
```

## Rate Limiting

No explicit rate limiting was observed, but standard courtesy applies - avoid excessive requests.
