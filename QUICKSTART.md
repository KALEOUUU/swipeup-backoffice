# 🚀 Quick Start Guide

Panduan cepat untuk menjalankan SwipeUp Backend.

## 📋 Prerequisites Checklist

- [ ] Go 1.21+ installed
- [ ] PostgreSQL 15+ installed & running
- [ ] Git installed

## ⚡ 5-Minute Setup

### 1. Clone & Install

```bash
git clone https://github.com/KALEOUUU/Telkom-UMKM-POS-APP.git
cd swipeup-be
go mod download
```

### 2. Database Setup

```bash
# Create database
psql -U postgres -c "CREATE DATABASE kantin_pos;"

# Run migrations
cat migrations/*.sql | psql -U postgres -d kantin_pos
```

### 3. Environment Configuration

Create `.env` file:

```env
SERVER_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=kantin_pos
DB_SSLMODE=disable
JWT_SECRET=change-this-secret-key
```

### 4. Run Server

```bash
go run cmd/server/main.go
```

Server running at `http://localhost:8080` 🎉

## 🧪 Test API

### Register User

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin001",
    "password": "admin123",
    "role": "admin_stan"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin001",
    "password": "admin123"
  }'
```

Copy token dari response untuk authenticated requests:

```bash
# Example authenticated request
curl http://localhost:8080/api/auth/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## 📚 Next Steps

- Read full documentation: [README.md](README.md)
- Explore API with Bruno: Open `docs/` folder in Bruno
- Check architecture diagram in README.md
- See refactoring notes: [REFACTORING_SUMMARY.md](REFACTORING_SUMMARY.md)

## 🐛 Troubleshooting

### Database Connection Error

```bash
# Check PostgreSQL is running
pg_isready

# Check connection
psql -U postgres -d kantin_pos -c "SELECT 1;"
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>

# Or change port in .env
SERVER_PORT=8081
```

### Module Download Issues

```bash
# Clear cache and retry
go clean -modcache
go mod download
```

## 📞 Need Help?

- Full README: [README.md](README.md)
- API Docs: Browse `docs/` folder
- Issues: [GitHub Issues](https://github.com/KALEOUUU/Telkom-UMKM-POS-APP/issues)

---

**Happy Coding! 🚀**
