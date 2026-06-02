# KODA B7 GIN
[![License: MIT](https://img.shields.io/badge/License-MIT-blue)](https://opensource.org/license/mit)
<br>
Project for gin exercise by Koda batch 7

## TECH USED
- [![Gin-Gonic](https://img.shields.io/badge/Gin_Gonic-v1.12.0-green?logo=Gin&logoColor=white)](https://gin-gonic.com/en/)
- [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-17.4-blue?logo=PostgreSQL&logoColor=white)](https://www.postgresql.org/)
- [![Redis](https://img.shields.io/badge/Redis-7.4-blue?logo=Redis&logoColor=white)](https://img.shields.io/snapcraft/l/:package)

## FEATURES
- Authentication (Login, Register, Logout)
- Enter PIN & Change PIN
- Dashboard (Balance, Income, Expense, Chart)
- Top Up dengan berbagai metode pembayaran
- Transfer dengan verifikasi PIN
- History Transaksi dengan pagination
- Edit Profile dengan upload foto
- Change password & Change pin


## USAGE INSTRUCTION
### Environment Setup
1. Create your environment on the root directory named
```.env```
```
DB_HOST={YOUR_DB_HOST}
DB_PORT={YOUR_DB_PORT}
DB_PASS={YOUR_DB_HOST}
DB_NAME={YOUR_DB_HOST}
DB_USER={YOUR_DB_HOST}
```

### Running the Application
1. Clone the App
```bash
$ git clone url
```
2. Install dependency
```bash
$ go mod download
```


## ROUTES
### Features
| Endpoint | Method | Description |
| ---------- | ------------| ------------ |
| /auth/login | POST | Login |
| /auth/register | POST | Register | 

### Documentation
For complete documentation, visit ``/swagger/index.html``

## Changelog
| Version | Description |
| ---------- | ------------|
| latest | add feature by [akmalrian](https://github.com/Akmalrian)|

## HOW TO CONTRIBUTE (optional)
- Fork this repository
- Create your changes
- Pull Request

## LICENSE
this project is licensed under the MIT License

<!-- ## CONTACTS
[email](mailto:) -->

## RELATED PROJECT
[Frontend](https://github.com/Akmalrian/E-Wallet-Frontend.git)