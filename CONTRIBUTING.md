# 🤝 Contributing to SwipeUp Backend

Thank you for your interest in contributing to SwipeUp Backend! This document provides guidelines and instructions for contributing.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Maintain professional communication

## Getting Started

1. **Fork the Repository**

```bash
# Click "Fork" button on GitHub
# Clone your fork
git clone https://github.com/YOUR_USERNAME/Telkom-UMKM-POS-APP.git
cd swipeup-be
```

2. **Set Up Upstream Remote**

```bash
git remote add upstream https://github.com/KALEOUUU/Telkom-UMKM-POS-APP.git
git fetch upstream
```

3. **Create Feature Branch**

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/bug-description
```

## Development Workflow

### 1. Make Your Changes

- Write clean, readable code
- Follow existing code patterns
- Add comments for complex logic
- Update documentation if needed

### 2. Test Your Changes

```bash
# Run tests
go test ./...

# Build to check compilation
go build ./cmd/server

# Manual testing
go run cmd/server/main.go
```

### 3. Update Documentation

- Update README.md if you add features
- Add Bruno API documentation in `docs/`
- Update CHANGELOG if exists
- Add comments to complex functions

## Coding Standards

### Go Code Style

```go
// ✅ Good
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
    var user models.User
    if err := s.db.First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

// ❌ Bad (no error handling)
func (s *UserService) GetUser(id uint) *models.User {
    var user models.User
    s.db.First(&user, id)
    return &user
}
```

### Naming Conventions

- **Functions**: CamelCase, start with verb (`GetUser`, `CreateOrder`)
- **Variables**: camelCase (`userId`, `orderTotal`)
- **Constants**: UPPER_CASE (`MAX_RETRY`, `DEFAULT_TIMEOUT`)
- **Interfaces**: End with `-er` (`Reader`, `Writer`, `Handler`)

### Project Structure

When adding new features:

```
internal/
├── models/          # Add your model here
│   └── your_model.go
├── services/        # Add business logic
│   └── your_service.go
├── handlers/        # Add HTTP handlers
│   └── your_handler.go
└── middleware/      # Add middleware if needed
    └── your_middleware.go
```

### Error Handling

```go
// ✅ Always handle errors
if err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
}

// ✅ Use helper functions
BadRequestResponse(c, "Invalid request", err)

// ❌ Don't ignore errors
user, _ := service.GetUser(id) // Bad!
```

### Database Operations

```go
// ✅ Use transactions for multiple operations
tx := s.db.Begin()
if err := tx.Create(&user).Error; err != nil {
    tx.Rollback()
    return err
}
if err := tx.Create(&profile).Error; err != nil {
    tx.Rollback()
    return err
}
tx.Commit()

// ✅ Use preload for relations
db.Preload("Menu").First(&stan, id)

// ❌ N+1 query problem
for _, stan := range stans {
    db.Find(&stan.Menu) // Bad!
}
```

### Response Format

Always use helper functions:

```go
// ✅ Use helpers
SuccessResponse(c, "User created", user)
BadRequestResponse(c, "Invalid input", err)
InternalErrorResponse(c, "Database error", err)

// ❌ Manual response
c.JSON(200, gin.H{"success": true, "data": user})
```

## Commit Guidelines

### Commit Message Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```bash
# Good commit messages
git commit -m "feat(auth): add JWT token refresh endpoint"
git commit -m "fix(cart): resolve checkout total calculation bug"
git commit -m "docs(readme): update installation instructions"
git commit -m "refactor(handlers): extract common pagination helper"

# Bad commit messages
git commit -m "update"
git commit -m "fix bug"
git commit -m "changes"
```

### Commit Best Practices

- One logical change per commit
- Write clear, descriptive messages
- Reference issue numbers: `fix(auth): resolve login error (#123)`
- Keep commits focused and small

## Pull Request Process

### 1. Update Your Branch

```bash
# Sync with upstream
git fetch upstream
git rebase upstream/main

# Or merge if preferred
git merge upstream/main
```

### 2. Push Your Changes

```bash
git push origin feature/your-feature-name
```

### 3. Create Pull Request

On GitHub:

1. Click "New Pull Request"
2. Select your branch
3. Fill in the template:

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] All tests pass
- [ ] Manual testing done
- [ ] Bruno collection updated

## Checklist
- [ ] Code follows project style
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
- [ ] Commits follow guidelines
```

### 4. Code Review Process

- Respond to feedback promptly
- Make requested changes
- Keep discussions professional
- Ask questions if unclear

### 5. After Approval

- Squash commits if requested
- Wait for maintainer to merge
- Delete your branch after merge

## Testing Guidelines

### Write Tests

```go
// Example test
func TestUserService_GetUserByID(t *testing.T) {
    // Setup
    db := setupTestDB()
    service := NewUserService(db)
    
    // Test
    user, err := service.GetUserByID(1)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "testuser", user.Username)
}
```

### Test Coverage

- Aim for >70% coverage
- Test happy paths
- Test error cases
- Test edge cases

## Documentation

### Code Comments

```go
// ✅ Good comment
// GetUserByID retrieves a user from database by ID.
// Returns error if user not found or database error occurs.
func GetUserByID(id uint) (*User, error) {
    // ...
}

// ❌ Redundant comment
// This function gets a user
func GetUserByID(id uint) (*User, error) {
    // ...
}
```

### API Documentation

Add Bruno documentation for new endpoints:

```
docs/
└── YourFeature/
    ├── folder.bru          # Feature overview
    ├── Create.bru          # POST endpoint
    ├── GetAll.bru          # GET list endpoint
    └── GetByID.bru         # GET single endpoint
```

## Questions?

- 📖 Check [README.md](README.md)
- 🚀 See [QUICKSTART.md](QUICKSTART.md)
- 💬 Open a [Discussion](https://github.com/KALEOUUU/Telkom-UMKM-POS-APP/discussions)
- 🐛 Report [Issues](https://github.com/KALEOUUU/Telkom-UMKM-POS-APP/issues)

---

**Thank you for contributing! 🎉**
