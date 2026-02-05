# Role-Based Architecture Documentation

## Overview

This document describes the refactored role-based architecture for the SwipeUp Backend. The system now provides clear separation of concerns for three distinct user roles:

1. **Superadmin** - Global administrator managing all stans and system-wide operations
2. **Stan Admin (admin_stan)** - Stan owner managing their own stan, menus, discounts, and payments
3. **Student (siswa)** - Customer making purchases and managing their profile

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                         API Gateway                              │
│                    (cmd/server/main.go)                          │
└─────────────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│  Superadmin   │    │  Stan Admin   │    │   Student     │
│   Routes      │    │   Routes      │    │   Routes      │
└───────────────┘    └───────────────┘    └───────────────┘
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│Superadmin     │    │StanAdmin      │    │Student        │
│Handler        │    │Handler        │    │Handler        │
└───────────────┘    └───────────────┘    └───────────────┘
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐    ┌───────────────┐    ┌───────────────┐
│Superadmin     │    │StanAdmin      │    │Student        │
│Service        │    │Service        │    │Service        │
└───────────────┘    └───────────────┘    └───────────────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              │
                              ▼
                    ┌───────────────────┐
                    │  Base Services   │
                    │  (User, Stan,   │
                    │   Menu, etc.)   │
                    └───────────────────┘
                              │
                              ▼
                    ┌───────────────────┐
                    │   Database       │
                    │   (PostgreSQL)   │
                    └───────────────────┘
```

## Role Definitions

### 1. Superadmin

**Purpose:** Global administrator managing all stans and system-wide operations

**Capabilities:**
- User management (CRUD operations)
- View all siswa profiles
- View all stans
- Create, update, delete global discounts
- View revenue reports for all stans
- View statistics for all stans
- View and manage activity logs
- System-level monitoring and analytics

**API Routes:**
```
POST   /api/auth/login
POST   /api/auth/register
POST   /api/auth/register-admin-stan
GET    /api/auth/profile

# Superadmin Only
POST   /api/superadmin/users
GET    /api/superadmin/users
GET    /api/superadmin/users/:id
PUT    /api/superadmin/users/:id
DELETE /api/superadmin/users/:id

GET    /api/superadmin/siswa
GET    /api/superadmin/siswa/:id
DELETE /api/superadmin/siswa/:id

GET    /api/superadmin/stan
GET    /api/superadmin/stan/:id
DELETE /api/superadmin/stan/:id

GET    /api/superadmin/revenue
GET    /api/superadmin/revenue/stan/:id
GET    /api/superadmin/revenue/all-stans

GET    /api/superadmin/statistics
GET    /api/superadmin/statistics/stan/:id

GET    /api/superadmin/discounts/global
POST   /api/superadmin/discounts/global
PUT    /api/superadmin/discounts/global/:id
DELETE /api/superadmin/discounts/global/:id

GET    /api/superadmin/activity-logs
GET    /api/superadmin/activity-logs/date-range
GET    /api/superadmin/activity-logs/stats
DELETE /api/superadmin/activity-logs/clean
```

### 2. Stan Admin (admin_stan)

**Purpose:** Stan owner managing their own stan, menus, discounts, and payments

**Capabilities:**
- Manage stan profile (name, owner info, photo)
- Manage payment settings (cash, QRIS)
- Create, update, delete menu items
- Manage menu inventory (stock adjustments)
- Create stan-level discounts
- Create menu-level discounts
- View and update transaction status
- View revenue for their stan
- View transactions for their stan

**API Routes:**
```
POST   /api/auth/login
POST   /api/auth/register-admin-stan
GET    /api/auth/profile

# Stan Admin Only
GET    /api/admin-stan/stan/profile
PUT    /api/admin-stan/stan/profile
PUT    /api/admin-stan/stan/payment-settings

GET    /api/admin-stan/menu
POST   /api/admin-stan/menu
PUT    /api/admin-stan/menu/:id
DELETE /api/admin-stan/menu/:id
PUT    /api/admin-stan/menu/:id/stock
PATCH  /api/admin-stan/menu/:id/adjust-stock

GET    /api/admin-stan/discounts
GET    /api/admin-stan/discounts/active
POST   /api/admin-stan/discounts/stan
POST   /api/admin-stan/discounts/menu
PUT    /api/admin-stan/discounts/:id
DELETE /api/admin-stan/discounts/:id

GET    /api/admin-stan/transactions
GET    /api/admin-stan/transactions/date-range
PUT    /api/admin-stan/transactions/:id/status

GET    /api/admin-stan/revenue
```

### 3. Student (siswa)

**Purpose:** Customer making purchases and managing their profile

**Capabilities:**
- Manage student profile (name, class, address, phone)
- Add items to cart
- View cart
- Update cart item quantity
- Remove items from cart
- Clear cart
- Checkout cart to create transaction
- View their transaction history
- View transaction details

**API Routes:**
```
POST   /api/auth/login
POST   /api/auth/register
GET    /api/auth/profile

# Student Only
GET    /api/student/profile
PUT    /api/student/profile

POST   /api/student/cart
GET    /api/student/cart
PUT    /api/student/cart/:id
DELETE /api/student/cart/:id
DELETE /api/student/cart/clear
POST   /api/student/cart/checkout

GET    /api/student/transactions
GET    /api/student/transactions/:id
```

### 4. Public Routes (All Authenticated Users)

**Purpose:** Public access to view stans, menus, and discounts

**API Routes:**
```
GET    /api/public/stan
GET    /api/public/stan/:id
GET    /api/public/menu
GET    /api/public/menu/:id
GET    /api/public/menu/by-stan
GET    /api/public/menu/available
GET    /api/public/menu/search
GET    /api/public/discounts/active-by-stan
```

## Middleware

### Authentication Middleware
- **`AuthMiddleware`**: Validates JWT tokens and sets user context
- **`OptionalAuthMiddleware`**: Allows requests without auth but sets user info if token provided

### Role-Based Authorization Middleware
- **`SuperAdminOnly`**: Restricts access to superadmin only
- **`AdminStanOnly`**: Restricts access to admin_stan only
- **`RoleMiddleware("siswa")`**: Restricts access to siswa only
- **`AdminAccess`**: Allows access to both superadmin and admin_stan

### Resource Ownership Middleware
- **`StanOwnerOnly`**: Verifies user owns the stan resource
- **`SiswaOwnerOnly`**: Verifies user owns the siswa profile
- **`MenuStanOwnerOnly`**: Verifies user owns the stan that menu belongs to
- **`TransaksiStanOwnerOnly`**: Verifies user owns the stan that transaction belongs to
- **`DiskonStanOwnerOnly`**: Verifies user owns the stan that discount belongs to

## Service Layer

### Base Services
- **`UserService`**: User CRUD operations
- **`SiswaService`**: Student profile management
- **`StanService`**: Stan profile management
- **`MenuService`**: Menu CRUD and inventory management
- **`DiskonService`**: Discount management
- **`TransaksiService`**: Transaction management
- **`CartService`**: Shopping cart management
- **`AuthService`**: Authentication and authorization
- **`ActivityLogService`**: Activity logging

### Role-Specific Services

#### SuperadminService
```go
type SuperadminService struct {
    db *gorm.DB
}

// Methods:
- GetRevenueByStanID(stanID, startDate, endDate) -> StanRevenue
- GetAllStanRevenue(startDate, endDate) -> []StanRevenue
- GetRevenueReport(startDate, endDate) -> RevenueReport
- GetGlobalDiscounts() -> []Diskon
- CreateGlobalDiscount(diskon) -> error
- UpdateGlobalDiscount(id, updates) -> error
- DeleteGlobalDiscount(id) -> error
- GetStanStatistics(stanID, startDate, endDate) -> StanStatistics
- GetAllStanStatistics(startDate, endDate) -> []StanStatistics
```

#### StanAdminService
```go
type StanAdminService struct {
    db            *gorm.DB
    stanService   *StanService
    menuService   *MenuService
    diskonService *DiskonService
}

// Methods:
- GetStanByUserID(userID) -> Stan
- UpdateStanProfile(userID, updates) -> error
- CreateMenu(userID, menu) -> error
- UpdateMenu(userID, menuID, updates) -> error
- DeleteMenu(userID, menuID) -> error
- GetMenusByStan(userID) -> []Menu
- CreateStanDiscount(userID, diskon) -> error
- CreateMenuDiscount(userID, diskon, menuIDs) -> error
- UpdateDiscount(userID, diskonID, updates) -> error
- DeleteDiscount(userID, diskonID) -> error
- GetDiscountsByStan(userID) -> []Diskon
- GetActiveDiscountsByStan(userID) -> []Diskon
- GetTransactionsByStan(userID) -> []Transaksi
- GetTransactionsByStanAndDateRange(userID, startDate, endDate) -> []Transaksi
- UpdateTransactionStatus(userID, transaksiID, status) -> error
- GetStanRevenue(userID, startDate, endDate) -> (float64, int, error)
- UpdatePaymentSettings(userID, acceptCash, acceptQris, qrisImage) -> error
```

#### StudentService
```go
type StudentService struct {
    db               *gorm.DB
    siswaService     *SiswaService
    cartService      *CartService
    transaksiService *TransaksiService
}

// Methods:
- GetSiswaByUserID(userID) -> Siswa
- UpdateSiswaProfile(userID, updates) -> error
- AddToCart(siswaID, cart) -> error
- GetCart(siswaID) -> ([]Cart, int, float64, error)
- UpdateCartItem(cartID, qty) -> error
- RemoveFromCart(cartID) -> error
- ClearCart(siswaID) -> error
- CheckoutCart(siswaID, stanID) -> (*Transaksi, []DetailTransaksi, error)
- GetTransactions(siswaID) -> []Transaksi
- GetTransactionByID(siswaID, transaksiID) -> (*Transaksi, error)
```

## Handler Layer

### Base Handlers
- **`UserHandler`**: User CRUD operations
- **`SiswaHandler`**: Student profile management
- **`StanHandler`**: Stan profile management
- **`MenuHandler`**: Menu CRUD and inventory management
- **`DiskonHandler`**: Discount management
- **`TransaksiHandler`**: Transaction management
- **`CartHandler`**: Shopping cart management
- **`AuthHandler`**: Authentication and authorization
- **`ActivityLogHandler`**: Activity logging

### Role-Specific Handlers

#### SuperadminHandler
```go
type SuperadminHandler struct {
    superadminService *services.SuperadminService
    stanService      *services.StanService
    diskonService    *services.DiskonService
}

// Methods:
- GetRevenueByStanID(c *gin.Context)
- GetAllStanRevenue(c *gin.Context)
- GetRevenueReport(c *gin.Context)
- GetGlobalDiscounts(c *gin.Context)
- CreateGlobalDiscount(c *gin.Context)
- UpdateGlobalDiscount(c *gin.Context)
- DeleteGlobalDiscount(c *gin.Context)
- GetStanStatistics(c *gin.Context)
- GetAllStanStatistics(c *gin.Context)
```

#### StanAdminHandler
```go
type StanAdminHandler struct {
    stanAdminService *services.StanAdminService
    menuService     *services.MenuService
}

// Methods:
- GetStanProfile(c *gin.Context)
- UpdateStanProfile(c *gin.Context)
- UpdatePaymentSettings(c *gin.Context)
- CreateMenu(c *gin.Context)
- UpdateMenu(c *gin.Context)
- DeleteMenu(c *gin.Context)
- GetMenus(c *gin.Context)
- UpdateStock(c *gin.Context)
- AdjustStock(c *gin.Context)
- CreateStanDiscount(c *gin.Context)
- CreateMenuDiscount(c *gin.Context)
- UpdateDiscount(c *gin.Context)
- DeleteDiscount(c *gin.Context)
- GetDiscounts(c *gin.Context)
- GetActiveDiscounts(c *gin.Context)
- GetTransactions(c *gin.Context)
- GetTransactionsByDateRange(c *gin.Context)
- UpdateTransactionStatus(c *gin.Context)
- GetRevenue(c *gin.Context)
```

#### StudentHandler
```go
type StudentHandler struct {
    studentService *services.StudentService
}

// Methods:
- GetSiswaProfile(c *gin.Context)
- UpdateSiswaProfile(c *gin.Context)
- AddToCart(c *gin.Context)
- GetCart(c *gin.Context)
- UpdateCartItem(c *gin.Context)
- RemoveFromCart(c *gin.Context)
- ClearCart(c *gin.Context)
- CheckoutCart(c *gin.Context)
- GetTransactions(c *gin.Context)
- GetTransactionByID(c *gin.Context)
```

## Data Models

### User Model
```go
type User struct {
    ID        uint      `json:"id"`
    Username  string    `json:"username"`
    Password  string    `json:"-"`
    Role      UserRole  `json:"role"` // superadmin, admin_stan, siswa
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Stan      *Stan     `json:"stan,omitempty"`
    Siswa     *Siswa    `json:"siswa,omitempty"`
}
```

### Stan Model
```go
type Stan struct {
    ID          uint   `json:"id"`
    NamaStan    string `json:"nama_stan"`
    NamaPemilik string `json:"nama_pemilik"`
    Telp        string `json:"telp"`
    Foto        string `json:"foto"`
    QrisImage   string `json:"qris_image"`
    AcceptCash  bool   `json:"accept_cash"`
    AcceptQris  bool   `json:"accept_qris"`
    IDUser      uint   `json:"id_user"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### Siswa Model
```go
type Siswa struct {
    ID        uint      `json:"id"`
    NamaSiswa string    `json:"nama_siswa"`
    Kelas     string    `json:"kelas"`
    Alamat    string    `json:"alamat"`
    Telp      string    `json:"telp"`
    IDUser    uint      `json:"id_user"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## Security Considerations

1. **JWT Token Validation**: All protected routes require a valid JWT token
2. **Role-Based Access Control**: Each role has specific permissions
3. **Resource Ownership**: Stan admins can only access their own resources
4. **Input Validation**: All inputs are validated using Gin binding
5. **SQL Injection Prevention**: Using GORM ORM with parameterized queries
6. **Password Hashing**: Using bcrypt for password hashing

## Testing Recommendations

### Superadmin Testing
1. Test user CRUD operations
2. Test global discount management
3. Test revenue reporting
4. Test stan statistics
5. Test activity log management

### Stan Admin Testing
1. Test stan profile management
2. Test menu CRUD operations
3. Test discount management
4. Test transaction status updates
5. Test revenue tracking
6. Test payment settings

### Student Testing
1. Test profile management
2. Test cart operations
3. Test checkout process
4. Test transaction viewing
5. Test permission restrictions

## Migration Guide

### For Existing Superadmin Users
- Use new `/api/superadmin/*` routes for superadmin-specific operations
- Access revenue reports via `/api/superadmin/revenue`
- Manage global discounts via `/api/superadmin/discounts/global`

### For Existing Stan Admin Users
- Use new `/api/admin-stan/*` routes for stan admin operations
- Manage stan profile via `/api/admin-stan/stan/profile`
- Manage menus via `/api/admin-stan/menu`
- Manage discounts via `/api/admin-stan/discounts`
- View transactions via `/api/admin-stan/transactions`

### For Existing Student Users
- Use new `/api/student/*` routes for student operations
- Manage profile via `/api/student/profile`
- Manage cart via `/api/student/cart`
- View transactions via `/api/student/transactions`

## Future Enhancements

1. **Notification System**: Notify stan admins of new orders
2. **Analytics Dashboard**: More detailed analytics for superadmin
3. **Multi-tenancy Support**: Support for multiple schools/organizations
4. **API Rate Limiting**: Prevent abuse of API endpoints
5. **Caching**: Implement Redis caching for frequently accessed data
6. **WebSocket Support**: Real-time updates for orders
7. **Mobile App API**: Dedicated API endpoints for mobile apps
