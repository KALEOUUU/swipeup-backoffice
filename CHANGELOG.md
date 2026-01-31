# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Complete documentation (README.md, QUICKSTART.md, CONTRIBUTING.md)
- .gitignore file for Go projects
- .env.example for environment configuration template

## [1.0.0] - 2026-01-31

### Added
- **Authentication System**: JWT-based authentication with bcrypt password hashing
- **Cart System**: Full shopping cart functionality for mobile app
  - Add to cart
  - Update cart items
  - Remove from cart
  - Clear cart
  - Checkout cart
- **Activity Logs**: User activity tracking and analytics
  - Track login, add_to_cart, checkout actions
  - Query by user, action, date range
  - Activity statistics
  - Automatic log cleanup
- **User Management**: Complete CRUD operations for users
- **Siswa (Student) Management**: Student profile management
- **Stan (Booth) Management**: Booth/stall management for vendors
- **Menu Management**: Food/drink menu CRUD with search functionality
- **Transaction Management**: Order processing and tracking
  - Create transactions with details
  - Update transaction status
  - Query by student or booth
- **Discount System**: Flexible discount management
  - Global and per-booth discounts
  - Percentage and fixed amount discounts
  - Assign/remove discounts to menu items
- **Role-Based Access Control**: Admin and student roles
- **Helper Functions**: Reusable utilities for DRY code
  - Response helpers (Success, Error, BadRequest)
  - Context helpers (GetUserIDFromContext)
  - Client info extraction (GetClientInfo)
  - Pagination helper (ParsePaginationParams)

### Changed
- Refactored activity log handler to use standard response helpers
- Improved code organization with base services and handlers
- Standardized error handling across all endpoints

### Removed
- Deprecated `product.go` and `transaction.go` handlers (6 files)
- Unused imports and duplicate functions
- Dead code (~200 lines)

### Fixed
- Code duplication issues (~360 lines reduced)
- Inconsistent response formats
- Missing authentication middleware imports

### Security
- JWT token authentication (24-hour expiry)
- Password hashing with bcrypt
- Protected routes with middleware
- Input validation and sanitization
- SQL injection prevention with GORM

## [0.2.0] - 2026-01-29

### Added
- Base service layer for generic CRUD operations
- Diskon (discount) system implementation
- Menu search functionality
- Transaction detail tracking

### Changed
- Database schema improvements with proper indexes
- Enhanced GORM models with relationships

## [0.1.0] - 2026-01-28

### Added
- Initial project setup
- Basic CRUD operations for core entities
- PostgreSQL database integration
- Gin web framework setup
- GORM ORM integration
- Database migrations

---

## Version History

### [1.0.0] - Major Release
**Focus**: Mobile app support with cart, authentication, and activity tracking

**Stats**:
- 9 main features
- 45+ API endpoints
- 8 database tables
- Complete Bruno API documentation
- Zero code duplication
- Clean, maintainable codebase

### [0.2.0] - Feature Enhancement
**Focus**: Discount system and search functionality

### [0.1.0] - Initial Release
**Focus**: Core POS functionality

---

## Migration Guide

### From 0.2.0 to 1.0.0

**Database Migrations**:
```bash
# Run new migrations
psql -U postgres -d kantin_pos -f migrations/0005_add_cart_table.sql
psql -U postgres -d kantin_pos -f migrations/0006_add_activity_logs_table.sql
```

**Environment Variables**:
```env
# Add to .env
JWT_SECRET=your-secret-key-here
```

**Code Changes**:
- All endpoints now require authentication (except login/register)
- Response format standardized across all endpoints
- Deprecated endpoints removed (product, transaction)

**Breaking Changes**:
- Product and Transaction models removed (use Menu and Transaksi)
- Authentication now required for all API calls
- Response format changed to consistent structure

---

## Upcoming Features

### v1.1.0 (Planned)
- [ ] WebSocket support for real-time updates
- [ ] File upload for menu images
- [ ] Payment gateway integration
- [ ] Email notifications

### v1.2.0 (Planned)
- [ ] Report generation (PDF)
- [ ] Analytics dashboard API
- [ ] Redis caching layer
- [ ] API rate limiting

### v2.0.0 (Future)
- [ ] GraphQL API
- [ ] Microservices architecture
- [ ] Message queue (RabbitMQ/Kafka)
- [ ] Advanced analytics

---

For more details on each version, see the [Git commit history](https://github.com/KALEOUUU/Telkom-UMKM-POS-APP/commits/main).
