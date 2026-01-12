# Multi-Database Support

Hướng dẫn cách chuyển đổi giữa các loại database trong ứng dụng.

## Supported Drivers

| Driver | Status | Package |
|--------|--------|---------|
| MySQL | ✅ Ready | `github.com/go-sql-driver/mysql` |
| PostgreSQL | ✅ Ready | `github.com/lib/pq` |

## Configuration

Chỉnh sửa file `.env` để chuyển đổi database:

### MySQL
```env
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=secret
DB_NAME=my_database
```

### PostgreSQL
```env
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=my_database
DB_SSLMODE=disable
```

## How It Works

1. Ứng dụng đọc `DB_DRIVER` từ `.env`
2. Factory function `storage.NewDatabase()` tạo connection tương ứng
3. Migrations tự động chạy khi server khởi động

## Adding New Database Driver

1. Tạo file adapter mới trong `internal/adapter/storage/`:
   ```go
   // newdb_adapter.go
   type NewDBDatabase struct {
       db  *sql.DB
       cfg *config.DatabaseConfig
   }
   
   func NewNewDBDatabase(cfg *config.DatabaseConfig, timezone string) (*NewDBDatabase, error) {
       // Connection logic
   }
   
   func (n *NewDBDatabase) GetDB() *sql.DB { return n.db }
   func (n *NewDBDatabase) Close() error { return n.db.Close() }
   func (n *NewDBDatabase) DriverName() string { return "newdb" }
   func (n *NewDBDatabase) RunMigrations() error { /* ... */ }
   ```

2. Thêm case vào factory trong `database.go`:
   ```go
   case "newdb":
       return NewNewDBDatabase(cfg, timezone)
   ```

3. Thêm driver vào `SupportedDrivers` trong `config/database.go`
