# API Documentation

## Base URL
`http://localhost:8080`

## Endpoints

### Health
- `GET /health` - Health check
- `GET /ready` - Readiness check
- `GET /api/v1/info` - App info

### Users
- `GET /api/v1/users` - List users
- `GET /api/v1/users/{id}` - Get user
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

### Items  
- `GET /api/v1/items` - List items
- `GET /api/v1/items/{id}` - Get item
- `POST /api/v1/items` - Create item
- `PUT /api/v1/items/{id}` - Update item
- `DELETE /api/v1/items/{id}` - Delete item
