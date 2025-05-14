# Product API

This project is a RESTful Product API developed using the Echo framework in Go.

## Requirements

- Go 1.21 or higher
- MongoDB
- Make (optional)
- Docker and Docker Compose (for containerized deployment)

## Installation

### Local Development

1. Clone the project:
```bash
git clone <repo-url>
cd product-api
```

2. Install dependencies:
```bash
go mod download
```

3. Start MongoDB:
```bash
# Make sure MongoDB is running
mongod
```

4. Start the application:
```bash
go run main.go
```

### Docker Deployment

1. Build and start the containers:
```bash
docker-compose up --build
```

2. To run in detached mode:
```bash
docker-compose up -d
```

3. To stop the containers:
```bash
docker-compose down
```

## API Endpoints

### List Products
```http
GET /api/products
```

### Get Product Details
```http
GET /api/products/:id
```

### Create Product
```http
POST /api/products
Content-Type: application/json

{
    "name": "Product Name",
    "description": "Product Description",
    "price": 99.99
}
```

### Update Product
```http
PUT /api/products/:id
Content-Type: application/json

{
    "name": "New Product Name",
    "description": "New Product Description",
    "price": 149.99
}
```

### Delete Product
```http
DELETE /api/products/:id
```

### Search Products
```http
GET /api/products/search?name=product&exact_match=false&min_price=10&max_price=100&sort_by_price=asc
```

## Project Structure

```
.
├── internal/
│   ├── controller/    # HTTP request handling layer
│   ├── service/       # Business logic layer
│   ├── repository/    # Database operations layer
│   ├── model/         # Data models
│   └── dto/          # Data transfer objects
├── main.go           # Main application file
├── Dockerfile        # Docker configuration
├── docker-compose.yml # Docker Compose configuration
└── README.md         # Project documentation
```

## Features

- RESTful API design
- MongoDB database integration
- Logging (Logrus)
- Error handling
- Product search and filtering
- Price-based sorting
- Docker support 