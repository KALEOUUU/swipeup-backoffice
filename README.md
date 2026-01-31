# SwipeUp Backend - Kantin POS System

> Backend API untuk sistem Point of Sale (POS) Kantin dengan arsitektur modern menggunakan Go, Gin, GORM, dan PostgreSQL.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-336791?style=flat&logo=postgresql)](https://www.postgresql.org/)
[![Gin Framework](https://img.shields.io/badge/Gin-Framework-00ADD8?style=flat)](https://gin-gonic.com/)

---

## 📋 Deskripsi Project

**SwipeUp** adalah sistem POS (Point of Sale) untuk kantin sekolah/kampus yang terdiri dari 3 komponen utama:

1. **Admin Web** - Dashboard untuk admin stan mengelola menu, transaksi, dan laporan
2. **Landing Page** - Website informasi untuk siswa dan pengunjung
3. **Mobile App** - Aplikasi mobile untuk siswa melakukan order, save to cart, dan track transaksi

**Backend ini** menyediakan REST API yang mendukung semua fitur dari ketiga platform tersebut.

### 🎯 Fitur Utama

- 🔐 **Authentication & Authorization** - JWT-based auth dengan role-based access (admin_stan, siswa)
- 🛒 **Cart System** - Save to cart untuk mobile shopping experience
- 📊 **Activity Logs** - Tracking user behavior dan analytics
- 💰 **Transaction Management** - Complete transaction flow dengan detail tracking
- 🎁 **Discount System** - Flexible discount dengan tipe global dan per-stan
- 📱 **Mobile-First API** - Optimized untuk mobile app performance
- 🔍 **Search & Filter** - Advanced menu search dan filtering
- 📈 **Real-time Updates** - Status transaksi real-time

---

## 🏗️ Arsitektur

```
swipeup-be/
├── cmd/
│   └── server/
│       └── main.go              # Entry point aplikasi
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── connection.go        # Database connection setup
│   ├── models/                  # Data models (GORM)
│   │   ├── user.go
│   │   ├── siswa.go
│   │   ├── stan.go
│   │   ├── menu.go
│   │   ├── transaksi.go
│   │   ├── detail_transaksi.go
│   │   ├── diskon.go
│   │   ├── menu_diskon.go
│   │   ├── cart.go
│   │   └── activity_log.go
│   ├── services/                # Business logic layer
│   │   ├── base_service.go      # Generic CRUD operations
│   │   ├── auth_service.go      # Authentication & JWT
│   │   ├── user_service.go
│   │   ├── siswa_service.go
│   │   ├── stan_service.go
│   │   ├── menu_service.go
│   │   ├── transaksi_service.go
│   │   ├── diskon_service.go
│   │   ├── cart_service.go
│   │   └── activity_log_service.go
│   ├── handlers/                # HTTP handlers (controllers)
│   │   ├── base_handler.go      # Reusable response helpers
│   │   ├── auth_handler.go
│   │   ├── user.go
│   │   ├── siswa_handler.go
│   │   ├── stan_handler.go
│   │   ├── menu_handler.go
│   │   ├── transaksi_handler.go
│   │   ├── diskon_handler.go
│   │   ├── cart_handler.go
│   │   └── activity_log_handler.go
│   └── middleware/
│       └── auth.go              # JWT authentication middleware
├── migrations/                  # Database migrations
│   ├── 0001_create_tables.sql
│   ├── 0002_kantin_schema.sql
│   ├── 0005_add_cart_table.sql
│   └── 0006_add_activity_logs_table.sql
├── docs/                        # API documentation (Bruno)
│   ├── 1-Users/
│   ├── 2-Siswa/
│   ├── 3-Stan/
│   ├── 4-Menu/
│   ├── 5-Transaksi/
│   ├── 6-Diskon/
│   ├── 7-Cart/
│   ├── 8-Auth/
│   └── 9-Activity-Logs/
├── pkg/
│   └── utils/
│       └── helpers.go
├── go.mod
├── go.sum
└── README.md
```

### 🎨 Design Patterns

- **Layered Architecture**: Separation of concerns (Models, Services, Handlers)
- **Repository Pattern**: Database abstraction dengan BaseService
- **DRY Principle**: Reusable helper functions dan base services
- **JWT Authentication**: Stateless authentication untuk scalability
- **Middleware Pattern**: Auth middleware untuk protected routes

### 🗄️ Database Schema

**8 Tabel Utama:**

```
users               siswa               stan                menu
├── id             ├── id              ├── id              ├── id
├── username       ├── id_user (FK)    ├── id_user (FK)    ├── id_stan (FK)
├── password       ├── nama_siswa      ├── nama_stan       ├── nama_menu
├── role           ├── alamat          ├── deskripsi       ├── harga
├── created_at     ├── telp            ├── lokasi          ├── deskripsi
└── updated_at     └── ...             └── ...             └── ...

transaksi           detail_transaksi    diskon              menu_diskon
├── id             ├── id              ├── id              ├── id
├── id_stan (FK)   ├── id_transaksi    ├── id_stan (FK)    ├── id_menu (FK)
├── id_siswa (FK)  ├── id_menu (FK)    ├── nama_diskon     ├── id_diskon (FK)
├── status         ├── qty             ├── tipe            └── ...
├── created_at     ├── harga           ├── nilai
└── ...            └── ...             └── ...

cart                activity_logs
├── id             ├── id
├── id_siswa (FK)  ├── id_user (FK)
├── id_menu (FK)   ├── action
├── qty            ├── description
├── created_at     ├── ip_address
└── updated_at     ├── user_agent
                   └── created_at
```

**Relasi:**
- User → Siswa (1:1)
- User → Stan (1:1)
- Stan → Menu (1:N)
- Stan → Transaksi (1:N)
- Siswa → Transaksi (1:N)
- Transaksi → DetailTransaksi (1:N)
- Menu → DetailTransaksi (1:N)
- Diskon ← MenuDiskon → Menu (M:N)
- Siswa → Cart (1:N)
- User → ActivityLog (1:N)

---

## 🚀 Getting Started

### Prerequisites

Pastikan Anda sudah menginstall:

- **Go** 1.21 atau lebih tinggi ([Download](https://go.dev/dl/))
- **PostgreSQL** 15+ ([Download](https://www.postgresql.org/download/))
- **Git** ([Download](https://git-scm.com/downloads))
- **Bruno** (optional, untuk testing API) ([Download](https://www.usebruno.com/))

### Installation

1. **Clone Repository**

```bash
git clone https://github.com/KALEOUUU/Telkom-UMKM-POS-APP.git
cd swipeup-be
```

2. **Install Dependencies**

```bash
go mod download
```

3. **Setup Database**

Buat database PostgreSQL:

```bash
# Login ke PostgreSQL
psql -U postgres

# Buat database
CREATE DATABASE kantin_pos;

# Keluar
\q
```

4. **Setup Environment Variables**

Buat file `.env` di root directory:

```env
# Server Configuration
SERVER_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=kantin_pos
DB_SSLMODE=disable

# JWT Secret (ganti dengan random string)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
```

5. **Run Database Migrations**

Jalankan migration files secara manual:

```bash
# Migration 1: Create basic tables
psql -U postgres -d kantin_pos -f migrations/0001_create_tables.sql

# Migration 2: Kantin schema
psql -U postgres -d kantin_pos -f migrations/0002_kantin_schema.sql

# Migration 3: Cart table
psql -U postgres -d kantin_pos -f migrations/0005_add_cart_table.sql

# Migration 4: Activity logs
psql -U postgres -d kantin_pos -f migrations/0006_add_activity_logs_table.sql
```

Atau jalankan semua sekaligus:

```bash
cat migrations/*.sql | psql -U postgres -d kantin_pos
```

6. **Build & Run**

```bash
# Build aplikasi
go build -o server cmd/server/main.go

# Jalankan server
./server
```

Atau langsung run tanpa build:

```bash
go run cmd/server/main.go
```

Server akan berjalan di `http://localhost:8080`

### 🧪 Testing API

Anda bisa test API menggunakan:

1. **Bruno** (Recommended)
   - Buka Bruno
   - Import folder `docs/` sebagai collection
   - Gunakan pre-configured requests

2. **cURL**

```bash
# Health check
curl http://localhost:8080/

# Register user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "siswa001",
    "password": "password123",
    "role": "siswa"
  }'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "siswa001",
    "password": "password123"
  }'
```

3. **Postman**
   - Import Bruno collection atau buat manual

---

## 📚 API Documentation

### Authentication

Semua endpoint (kecuali login/register) memerlukan JWT token di header:

```
Authorization: Bearer <your-jwt-token>
```

### Available Endpoints

#### 🔐 Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login dan dapatkan JWT token
- `GET /api/auth/profile` - Get user profile (authenticated)

#### 👤 Users
- `POST /api/users` - Create user
- `GET /api/users` - Get all users
- `GET /api/users/:id` - Get user by ID

#### 🎓 Siswa (Students)
- `POST /api/siswa` - Create siswa
- `GET /api/siswa` - Get all siswa
- `GET /api/siswa/:id` - Get siswa by ID
- `PUT /api/siswa/:id` - Update siswa
- `DELETE /api/siswa/:id` - Delete siswa
- `GET /api/siswa/by-user?user_id=1` - Get siswa by user ID

#### 🏪 Stan (Booth/Stall)
- `POST /api/stan` - Create stan
- `GET /api/stan` - Get all stan
- `GET /api/stan/:id` - Get stan by ID (with menu)
- `PUT /api/stan/:id` - Update stan
- `DELETE /api/stan/:id` - Delete stan
- `GET /api/stan/by-user?user_id=1` - Get stan by user ID

#### 🍔 Menu
- `POST /api/menu` - Create menu
- `GET /api/menu` - Get all menu
- `GET /api/menu/:id` - Get menu by ID
- `PUT /api/menu/:id` - Update menu
- `DELETE /api/menu/:id` - Delete menu
- `GET /api/menu/by-stan?stan_id=1` - Get menu by stan ID
- `GET /api/menu/search?name=nasi` - Search menu by name

#### 💳 Transaksi (Transactions)
- `POST /api/transaksi` - Create transaction
- `GET /api/transaksi` - Get all transactions
- `GET /api/transaksi/:id` - Get transaction by ID (with details)
- `PUT /api/transaksi/:id/status` - Update transaction status
- `GET /api/transaksi/by-siswa?siswa_id=1` - Get transactions by siswa
- `GET /api/transaksi/by-stan?stan_id=1` - Get transactions by stan

#### 🎁 Diskon (Discounts)
- `POST /api/diskon` - Create discount
- `GET /api/diskon` - Get all discounts
- `GET /api/diskon/active` - Get active discounts
- `GET /api/diskon/global` - Get global discounts
- `GET /api/diskon/:id` - Get discount by ID
- `PUT /api/diskon/:id` - Update discount
- `DELETE /api/diskon/:id` - Delete discount
- `POST /api/diskon/:id/assign` - Assign discount to menu
- `DELETE /api/diskon/:id/remove?menu_id=1` - Remove discount from menu

#### 🛒 Cart (Shopping Cart)
- `POST /api/cart` - Add item to cart
- `GET /api/cart?siswa_id=1` - Get cart items
- `PUT /api/cart/:id` - Update cart item quantity
- `DELETE /api/cart/:id` - Remove item from cart
- `DELETE /api/cart/clear?siswa_id=1` - Clear all cart items
- `POST /api/cart/checkout` - Checkout cart

#### 📊 Activity Logs
- `GET /api/activity-logs/user?user_id=1` - Get user activities
- `GET /api/activity-logs?action=login` - Get all activities (with filter)
- `GET /api/activity-logs/date-range?start_date=2024-01-01&end_date=2024-01-31` - Get activities by date range
- `GET /api/activity-logs/stats` - Get activity statistics
- `DELETE /api/activity-logs/clean?days=90` - Clean old logs

**Detail lengkap ada di folder `docs/` (Bruno collections)**

---

## 🔧 Development

### Project Structure Explained

```
internal/
├── models/          # Database models dengan GORM tags
├── services/        # Business logic & database operations
├── handlers/        # HTTP request handlers (controllers)
├── middleware/      # Middleware (auth, logging, dll)
├── config/          # Configuration management
└── database/        # Database connection setup
```

### Adding New Feature

1. **Create Model** di `internal/models/`
2. **Create Migration** di `migrations/`
3. **Create Service** di `internal/services/`
4. **Create Handler** di `internal/handlers/`
5. **Register Routes** di `cmd/server/main.go`
6. **Add Documentation** di `docs/`

### Code Style

- Follow Go conventions dan idioms
- Use `gofmt` untuk formatting
- Implement DRY principles
- Write meaningful comments
- Use consistent naming conventions

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/services/...
```

---

## 🛡️ Security

- ✅ **Password Hashing**: bcrypt untuk password security
- ✅ **JWT Authentication**: Secure token-based auth
- ✅ **SQL Injection Prevention**: GORM parameterized queries
- ✅ **CORS**: Configurable CORS headers
- ✅ **Input Validation**: Request body validation
- ✅ **Role-Based Access**: Admin vs Student permissions

### Environment Variables Security

**JANGAN** commit `.env` file ke Git! File sudah masuk `.gitignore`.

Untuk production, gunakan environment variables yang secure atau secret management service.

---

## 📦 Dependencies

### Main Dependencies

```go
github.com/gin-gonic/gin           // HTTP web framework
github.com/golang-jwt/jwt/v5       // JWT authentication
golang.org/x/crypto/bcrypt         // Password hashing
gorm.io/gorm                       // ORM untuk database
gorm.io/driver/postgres            // PostgreSQL driver
github.com/joho/godotenv          // Environment variables
```

### Development Tools

- **Air** - Live reload untuk development
- **Bruno** - API testing dan documentation

---

## 🚀 Deployment

### Build untuk Production

```bash
# Build dengan optimizations
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go

# Atau dengan version info
go build -ldflags="-X main.Version=1.0.0" -o server cmd/server/main.go
```

### Docker (Optional)

Buat `Dockerfile`:

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./server"]
```

Build & run:

```bash
docker build -t swipeup-backend .
docker run -p 8080:8080 --env-file .env swipeup-backend
```

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License.

---

## 👥 Team

**SwipeUp Backend Team**
- Backend Developer: [Your Name]
- Project: Telkom UMKM POS Application

---

## 📞 Support

Untuk pertanyaan atau issue:
- 📧 Email: your-email@example.com
- 🐛 Issues: [GitHub Issues](https://github.com/KALEOUUU/Telkom-UMKM-POS-APP/issues)

---

## 🎯 Roadmap

### ✅ Completed
- [x] Authentication & Authorization (JWT)
- [x] Cart System
- [x] Activity Logs
- [x] Complete CRUD operations
- [x] Discount system
- [x] Transaction management

### 🚧 In Progress
- [ ] Unit tests coverage
- [ ] Integration tests
- [ ] API rate limiting

### 📋 Planned
- [ ] WebSocket untuk real-time updates
- [ ] File upload untuk menu images
- [ ] Payment gateway integration
- [ ] Email notifications
- [ ] Report generation (PDF)
- [ ] Analytics dashboard
- [ ] Caching layer (Redis)

---

<div align="center">

**Made with ❤️ using Go & Gin**

⭐ Star this repository if you find it helpful!

</div>
