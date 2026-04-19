# iiceekiingfx.com - Forex Trading Analytics Platform

A comprehensive Forex trading ecosystem combining education, analytics, and premium services.

## 🏗️ Architecture

### Tech Stack
- **Frontend**: Next.js (App Router), React, TypeScript, TailwindCSS
- **Backend**: Go (Fiber framework), PostgreSQL
- **Worker**: C# (.NET) for MT5 synchronization
- **Cache**: Redis (optional)

### Features
- 🔐 JWT Authentication (access + refresh tokens)
- 📊 Portfolio Analytics (MyFXBook-style)
- 📈 Real-time Trading Metrics
- 🎓 Course & Mentorship System
- 📡 Trading Signals
- 📝 Trade Journal
- 🧮 Position Size Calculator
- 💼 Account Management Services

## 🚀 Quick Start

### Prerequisites
- Go 1.22+
- PostgreSQL 14+
- Node.js 18+ (for frontend)

### Backend Setup

1. **Clone and setup**
```bash
git clone <repository>
cd iiceekiingfx-website
```

2. **Install dependencies**
```bash
go mod download
```

3. **Setup environment**
```bash
cp .env.example .env
# Edit .env with your database credentials
```

4. **Setup database**
```bash
# Create database
createdb iiceekiingfx

# Run migrations
psql -d iiceekiingfx -f database/migrations/001_create_users_table.sql
psql -d iiceekiingfx -f database/migrations/002_create_trading_accounts_table.sql
psql -d iiceekiingfx -f database/migrations/003_create_trades_table.sql
psql -d iiceekiingfx -f database/migrations/004_create_courses_tables.sql
psql -d iiceekiingfx -f database/migrations/005_create_journal_and_signals_tables.sql
```

5. **Run the server**
```bash
go run cmd/server/main.go
```

Server will start on `http://localhost:8080`

### API Endpoints

#### Authentication
- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login
- `GET /api/auth/me` - Get current user (protected)

#### Dashboard
- `GET /api/dashboard/overview` - Dashboard overview
- `GET /api/dashboard/equity-curve` - Equity curve data
- `GET /api/dashboard/activity` - Activity feed

#### Portfolio
- `POST /api/portfolio/connect` - Connect MT5 account
- `GET /api/portfolio/accounts` - Get trading accounts
- `GET /api/portfolio/history` - Trading history

#### Trade Journal
- `GET /api/journal/` - Get journal entries
- `POST /api/journal/` - Create journal entry
- `PUT /api/journal/:id` - Update journal entry
- `DELETE /api/journal/:id` - Delete journal entry

#### Courses
- `GET /api/courses/` - Get courses
- `GET /api/courses/:id` - Get course details

#### Signals
- `GET /api/signals/` - Get signals
- `POST /api/signals/` - Create signal (admin)

## 📁 Project Structure

```
iiceekiingfx-website/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/                # HTTP handlers
│   ├── services/                # Business logic
│   ├── repositories/            # Data access layer
│   ├── models/                  # Data models
│   └── middleware/              # HTTP middleware
├── config/                      # Configuration
├── database/
│   └── migrations/              # SQL migrations
├── pkg/                         # Shared packages
└── .env.example                 # Environment variables template
```

## 🔐 Security

- JWT tokens for authentication
- Password hashing with bcrypt
- Encrypted trading credentials
- Input validation
- CORS protection
- Rate limiting (recommended)

## 🛠️ Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o bin/server cmd/server/main.go
```

## 📊 Database Schema

The platform uses PostgreSQL with the following main tables:
- `users` - User accounts and authentication
- `trading_accounts` - Connected MT5 accounts
- `trades` - Historical trade data
- `courses` - Educational content
- `lessons` - Course lessons
- `course_progress` - User progress tracking
- `trade_journal` - Manual trade journaling
- `signals` - Trading signals

## 🚀 Deployment

### Docker (Recommended)
```dockerfile
# Build and run with Docker Compose
docker-compose up -d
```

### Manual Deployment
1. Set up PostgreSQL database
2. Run migrations
3. Build and deploy the Go binary
4. Configure reverse proxy (nginx)
5. Set up SSL certificates

## 📝 TODO

- [ ] Redis caching implementation
- [ ] C# MT5 worker service
- [ ] Next.js frontend
- [ ] Advanced analytics calculations
- [ ] Email notifications
- [ ] Payment integration
- [ ] Admin dashboard

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
