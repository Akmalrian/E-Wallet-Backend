# E-Wallet App - Backend

[![License: MIT](https://img.shields.io/badge/License-MIT-blue)](https://opensource.org/license/mit)
<br>
Backend project for E-Wallet by Akmal Oktarian (Koda b-7 Fullstack web developer)

## FEATURES
- Authentication (Login, Register, Logout)
- Enter PIN & Change PIN
- Dashboard (Balance, Income, Expense, Chart)
- Top Up dengan berbagai metode pembayaran
- Transfer dengan verifikasi PIN
- History Transaksi dengan pagination
- Edit Profile dengan upload foto
- Change password & Change pin

## Technologies Used

- [![Go](https://img.shields.io/badge/Go-1.24-00ADD8?logo=go&logoColor=white)](https://go.dev/)
- [![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?logo=go&logoColor=white)](https://gin-gonic.com/)
- [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16.13-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
- [![Redis](https://img.shields.io/badge/Redis-8.6.3-FF4438?logo=redis&logoColor=white)](https://redis.io/)
- [![JWT](https://img.shields.io/badge/JWT-Auth-000000?logo=jsonwebtokens&logoColor=white)](https://jwt.io/)
- [![Swagger](https://img.shields.io/badge/Swagger-Docs-85EA2D?logo=swagger&logoColor=white)](https://swagger.io/)
- [![Docker](https://img.shields.io/badge/Docker-29.5.2-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)

## Environment

```bash
APP_HOST={your_app_host}
APP_PORT={your_app_port}

DB_HOST={your_database_host}
DB_PORT={your_database_port}
DB_USER={your_database_user}
DB_PASS={your_database_password}
DB_NAME={your_database_name}

JWT_ISSUER={your_jwt_issuer}
JWT_SECRET={your_jwt_secret}

RDB_USER={your_redis_user}
RDB_PASS={your_redis_password}

SMTP_PORT={your_smpt_port}
SMTP_USER={your_smpt_uset}
SMTP_PASSWORD={your_smpt_password}
SMTP_FROM_EMAIL={your_smpt_email}
```

## ⚙️ Installation

1. Clone the project

```sh
$ https://github.com/Akmalrian/E-Wallet-Backend.git
```

2. Navigate to project directory

```sh
$ cd Tickitz-backend
```

3. Install dependencies

```sh
$ go mod tidy
```

4. Setup your [environment](##-environment)

5. Install [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate#installation) for DB migration

6. Do the DB Migration

```sh
$ migrate -database YOUR_DATABASE_URL -path ./db/migrations up
```

or if you install Makefile run command

```sh
$ make migrate-up
```

and seeding data

```sh
$ make db-seed
```

7. Run the project

```sh
$ go run ./cmd/main.go
```

## API Endpoint

| Method | Endpoint                  | Description                                                             |
| ------ | --------------------------| ------------------------------------------------------------------------|
| POST   | `/auth/register`          | Register new user                                                       |
| POST   | `/auth/login`             | Login user                                                              |
| DELETE | `/auth/logout`            | Logout                                                                  |
| POST   | `/auth/enter-pin`         | Request an enter pin                                                    |
| GET    | `/users/dashboard`        | Get for dashboard detail                                                |
| GET    | `/users/receiver`         | Get a new recever for transaction                                       |
| GET    | `/users/profile`          | Get user profile                                                        |
| PATCH  | `/users/password`         | Change password user                                                    |
| PATCH  | `/users/pin`              | Change pin user                                                         |
| POST   | `/transactions/topup`     | Update user data profile & change password                              |
| POST   | `/transactions/transfer`  | Get list of Order History for logged-in User                            |
| GET    | `/transactions/history`   | Get detailed information for specific order history by user logged-in   |

Full interactive docs available at `/swagger/index.html` after running the server.

## Changelog
| Version | Description |
| ------- | ----------- |
| 1.0  | Setup Docker multi-stage build and docker-compose orchestration with PostgreSQL & Redis by [akmalrian](https://github.com/Akmalrian) |

## How to Contribute
- Fork this repository
- Create your changes
- Commit your changes (Please strictly follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) standard: `feat:`, `fix:`, `chore:`, `docs:`)
- Push to the branch
- Open a Pull Request

## LICENSE
this project is licensed under the MIT License

## RELATED PROJECT
[Frontend E-Wallet](https://github.com/Akmalrian/E-Wallet-Frontend.git)