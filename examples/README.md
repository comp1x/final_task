# CUSTOMER

### CreateOffice 

1. Открываем консоль/bash консоль.
2. Чтобы создать офис - прокидываем в консоль, не забываем передать свои значения:

```
curl -X POST -H "Content-Type: application/json" -d '{
  "name": "string",
  "address": "string"
}' "http://localhost:13101/customer/offices"
```

### GetOffices

```
curl "http://localhost:13101/customer/offices"
```

### CreateUser

```
curl -X POST -H "Content-Type: application/json" -d '{
    "name": "string",
    "office_uuid": "string"
}' "http://localhost:13101/customer/users"
```

### GetUsers

```
curl "http://localhost:13101/customer/users?office_uuid=uuid"
```

### GetActualMenu

```
curl "http://localhost:13101/customer/users/orders"
```

### CreateOrder

```
curl -X POST -H "Content-Type: application/json" -d '{
  "user_uuid": "string",
  "salads": [
    {
      "count": 0,
      "product_uuid": "string"
    }
  ],
  "garnishes": [
    {
      "count": 0,
      "product_uuid": "string"
    }
  ],
  "meats": [
    {
      "count": 0,
      "product_uuid": "string"
    }
  ],
  "soups": [
    {
      "count": 0,
      "product_uuid": "string"
    }
  ],
  "drinks": [
    {
      "count": 0,
      "product_uuid": "string"
    }
  ],
  "desserts": [
    {
      "count": 0,
      "product_uuid": "string"
    }
  ]
}' "http://localhost:13101/customer/users/orders"
```

# RESTAURANT

### GetMenu

```
curl "http://localhost:13103/restaurant/menu?on_date=2023-06-05T15:04:05.999999999Z"
```

### CreateMenu

```
curl -X POST -H "Content-Type: application/json" -d '{
  "on_date": "2023-06-06T16:14:43.401Z",
  "opening_record_at": "2023-06-06T10:00:00.401Z",
  "closing_record_at": "2023-06-06T22:00:00.401Z",
  "salads": [
    "uuid"
  ],
  "garnishes": [
    "uuid"
  ],
  "meats": [
    "uuid"
  ],
  "soups": [
    "uuid"
  ],
  "drinks": [
    "uuid"
  ],
  "desserts": [
    "uuid"
  ]
}' "http://localhost:13103/restaurant/menu"
```

### GetUpToDateOrderList

```
curl "http://localhost:13103/restaurant/orders"
```

### GetProductList

```
curl "http://localhost:13103/restaurant/products"
```

### CreateProduct

```
curl -X POST -H "Content-Type: application/json" -d '{
  "name": "string",
  "description": "string",
  "type": "PRODUCT_TYPE_UNSPECIFIED",
  "weight": 0,
  "price": 0
}' http://localhost:13103/restaurant/products
```

# STATISTICS

### GetAmountOfProfit

```
curl "http://localhost:13105/statistics/amount-of-profit?start_date=2023-06-06T12:14:43.401Z&end_date=2023-06-06T23:00:00.401Z"
```

### TopProducts

```
curl "http://localhost:13105/statistics/top-products"
```
