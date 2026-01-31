# Kantin POS API Documentation

API untuk sistem Point of Sale (POS) Kantin dengan fitur lengkap untuk manajemen stan, menu, siswa, transaksi, dan diskon.

## Database Schema

Database terdiri dari 8 tabel utama:
- **users**: Menyimpan data pengguna (admin_stan, siswa)
- **stan**: Menyimpan data stan/booth kantin
- **siswa**: Menyimpan data siswa
- **menu**: Menyimpan data menu makanan/minuman
- **transaksi**: Menyimpan data transaksi
- **detail_transaksi**: Menyimpan detail item dalam transaksi
- **diskon**: Menyimpan data diskon
- **menu_diskon**: Junction table untuk relasi menu dan diskon

## Features

### ✨ DRY Code Architecture
- **BaseService**: Generic service layer untuk operasi CRUD umum
- **BaseHandler**: Reusable response handlers dan utility functions
- **Type-safe Models**: Menggunakan constants untuk enum values
- **Consistent Error Handling**: Standardized error responses

### 🔥 Advanced Features
- Soft delete (GORM DeletedAt)
- Eager loading dengan Preload
- Transaction support untuk operasi kompleks
- Query optimization dengan indexes
- Cascading deletes untuk referential integrity

## API Endpoints

### Users
- `POST /api/users` - Create user
- `GET /api/users` - Get all users
- `GET /api/users/:id` - Get user by ID

### Siswa (Students)
- `POST /api/siswa` - Create siswa
- `GET /api/siswa` - Get all siswa
- `GET /api/siswa/:id` - Get siswa by ID
- `PUT /api/siswa/:id` - Update siswa
- `DELETE /api/siswa/:id` - Delete siswa
- `GET /api/siswa/by-user?user_id=1` - Get siswa by user ID

### Stan (Booth/Stall)
- `POST /api/stan` - Create stan
- `GET /api/stan` - Get all stan
- `GET /api/stan/:id` - Get stan by ID (with menu)
- `PUT /api/stan/:id` - Update stan
- `DELETE /api/stan/:id` - Delete stan
- `GET /api/stan/by-user?user_id=1` - Get stan by user ID

### Menu
- `POST /api/menu` - Create menu
- `GET /api/menu` - Get all menu
- `GET /api/menu/:id` - Get menu by ID (with diskon)
- `PUT /api/menu/:id` - Update menu
- `DELETE /api/menu/:id` - Delete menu
- `GET /api/menu/by-stan?stan_id=1` - Get menu by stan
- `GET /api/menu/search?name=nasi` - Search menu by name

### Transaksi (Transactions)
- `POST /api/transaksi` - Create transaction with details
- `GET /api/transaksi` - Get all transactions
- `GET /api/transaksi/:id` - Get transaction by ID
- `PUT /api/transaksi/:id/status` - Update transaction status
- `GET /api/transaksi/by-siswa?siswa_id=1` - Get transactions by student
- `GET /api/transaksi/by-stan?stan_id=1` - Get transactions by stan

### Diskon (Discounts)
- `POST /api/diskon` - Create diskon
- `GET /api/diskon` - Get all diskon
- `GET /api/diskon/active` - Get active diskon
- `GET /api/diskon/:id` - Get diskon by ID
- `PUT /api/diskon/:id` - Update diskon
- `DELETE /api/diskon/:id` - Delete diskon
- `POST /api/diskon/:id/assign` - Assign diskon to menu
- `DELETE /api/diskon/:id/remove?menu_id=1` - Remove diskon from menu

## Request/Response Examples

### Create Transaction
```json
POST /api/transaksi
{
  "id_stan": 1,
  "id_siswa": 1,
  "details": [
    {
      "id_menu": 1,
      "qty": 2,
      "harga_beli": 15000
    },
    {
      "id_menu": 2,
      "qty": 1,
      "harga_beli": 10000
    }
  ]
}
```

### Update Transaction Status
```json
PUT /api/transaksi/1/status
{
  "status": "dimasak"
}
```

Status values: `belum dikonfirm`, `dimasak`, `diantar`, `sampai`

### Create Menu
```json
POST /api/menu
{
  "nama_makanan": "Nasi Goreng",
  "harga": 15000,
  "jenis": "makanan",
  "foto": "nasi-goreng.jpg",
  "deskripsi": "Nasi goreng spesial dengan telur",
  "id_stan": 1
}
```

## Running the Application

```bash
# Install dependencies
go mod download

# Run the server
go run cmd/server/main.go
```

## Code Structure

```
internal/
├── models/          # Data models dengan relations
│   ├── user.go
│   ├── siswa.go
│   ├── stan.go
│   ├── menu.go
│   ├── transaksi.go
│   ├── detail_transaksi.go
│   ├── diskon.go
│   └── menu_diskon.go
├── services/        # Business logic
│   ├── base_service.go      # Generic CRUD operations
│   ├── siswa_service.go
│   ├── stan_service.go
│   ├── menu_service.go
│   ├── transaksi_service.go
│   └── diskon_service.go
└── handlers/        # HTTP handlers
    ├── base_handler.go      # Response utilities
    ├── siswa_handler.go
    ├── stan_handler.go
    ├── menu_handler.go
    ├── transaksi_handler.go
    └── diskon_handler.go
```

## DRY Principles Applied

1. **BaseService**: Generic service untuk semua operasi CRUD
2. **BaseHandler**: Utility functions untuk responses
3. **Constants**: Type-safe enum values
4. **Preload Helper**: Consistent eager loading
5. **Error Handling**: Standardized error responses
6. **Transaction Support**: Reusable transaction patterns

## Technologies Used

- **Gin**: Web framework
- **GORM**: ORM
- **PostgreSQL**: Database
- **Go Generics**: Type-safe base service
