# Cashier App - Backend API

Backend API untuk aplikasi kasir/POS yang dibangun menggunakan Go dengan arsitektur layered.

## Tech Stack

- **Go 1.23+** - Programming language
- **Gin** - Web framework
- **GORM** - ORM untuk database
- **PostgreSQL** - Database
- **JWT** - Authentication
- **bcrypt** - Password hashing

## Struktur Project

```
backend/
├── config/          # Konfigurasi database dan migrasi
├── controllers/     # HTTP handlers
├── dto/             # Data Transfer Objects
├── middleware/      # Middleware (auth, cors)
├── models/          # Entity models
├── repositories/    # Data access layer
├── routes/          # Route definitions
├── services/        # Business logic layer
├── main.go          # Entry point
├── go.mod           # Dependencies
└── .env.example     # Environment variables template
```

## Arsitektur

Aplikasi ini menggunakan **Layered Architecture**:

1. **Controllers** - Menerima HTTP request dan mengembalikan response
2. **Services** - Business logic dan validasi
3. **Repositories** - Akses ke database
4. **Models** - Entity/table definitions
5. **DTOs** - Request/Response data structures

## Setup & Instalasi

### Prerequisites

- Go 1.23+
- PostgreSQL 14+

### 1. Clone dan masuk ke direktori

```bash
cd backend
```

### 2. Copy environment file

```bash
cp .env.example .env
```

### 3. Edit file .env sesuai konfigurasi database

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=cashier_db
DB_SSLMODE=disable
JWT_SECRET=your-super-secret-jwt-key
PORT=8080
```

### 4. Buat database

```sql
CREATE DATABASE cashier_db;
```

### 5. Install dependencies dan jalankan

```bash
go mod tidy
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/auth/login | Login user |
| POST | /api/v1/auth/register | Register user baru |

### Users (Admin only)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/users | Get semua users |
| GET | /api/v1/users/:id | Get user by ID |
| GET | /api/v1/users/profile | Get current user profile |
| PUT | /api/v1/users/:id | Update user |
| DELETE | /api/v1/users/:id | Delete user |

### Categories

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/categories | Get semua categories |
| GET | /api/v1/categories/:id | Get category by ID |
| POST | /api/v1/categories | Create category (Manager+) |
| PUT | /api/v1/categories/:id | Update category (Manager+) |
| DELETE | /api/v1/categories/:id | Delete category (Admin) |

### Products

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/products | Get semua products |
| GET | /api/v1/products/:id | Get product by ID |
| GET | /api/v1/products/category/:id | Get products by category |
| POST | /api/v1/products | Create product (Manager+) |
| PUT | /api/v1/products/:id | Update product (Manager+) |
| PATCH | /api/v1/products/:id/stock | Update stock (Manager+) |
| DELETE | /api/v1/products/:id | Delete product (Admin) |

### Transactions

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/transactions | Get semua transactions |
| GET | /api/v1/transactions/:id | Get transaction by ID |
| GET | /api/v1/transactions/code/:code | Get transaction by code |
| GET | /api/v1/transactions/today | Get today's transactions |
| GET | /api/v1/transactions/user/:id | Get transactions by user |
| POST | /api/v1/transactions | Create transaction |
| POST | /api/v1/transactions/:id/cancel | Cancel transaction (Manager+) |

### Settings (Admin only)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/settings | Get semua settings |
| GET | /api/v1/settings/:key | Get setting by key |
| GET | /api/v1/settings/store | Get store settings |
| GET | /api/v1/settings/payment | Get payment settings |
| PUT | /api/v1/settings | Update setting |
| PUT | /api/v1/settings/batch | Update multiple settings |

### Reports (Protected)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/v1/reports/dashboard | Get dashboard data |
| GET | /api/v1/reports/revenue/daily | Get daily revenue |
| GET | /api/v1/reports/revenue/range | Get revenue by date range |
| GET | /api/v1/reports/payment-distribution | Get payment distribution |
| GET | /api/v1/reports/products/top | Get top selling products |
| GET | /api/v1/reports/summary/monthly | Get monthly summary |
| GET | /api/v1/reports/export/transactions | Export transactions (Manager+) |

## Authentication

Semua endpoint (kecuali login dan register) memerlukan JWT token di header:

```
Authorization: Bearer <token>
```

## User Roles

- **admin** - Full access
- **manager** - CRUD products, categories, cancel transactions
- **cashier** - Create transactions, view data

## Default Admin Account

Setelah migrasi pertama, admin default akan dibuat:

- Email: `admin@cashier.com`
- Password: `admin123`

**PENTING:** Ubah password default setelah instalasi!

## Response Format

Semua response menggunakan format standar:

```json
{
  "status": true,
  "message": "Success message",
  "data": { ... }
}
```

Error response:

```json
{
  "status": false,
  "message": "Error message",
  "error": "Error details"
}
```

## Development

### Run dengan hot reload

```bash
# Install air
go install github.com/air-verse/air@latest

# Run
air
```

### Build

```bash
go build -o cashier-api main.go
```

## License

MIT
