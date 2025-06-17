# bzCommerce

[![Backend CI](https://github.com/bzelaznicki/bzCommerce/actions/workflows/backend-ci.yml/badge.svg)](https://github.com/bzelaznicki/bzCommerce/actions/workflows/backend-ci.yml) [![Frontend CI](https://github.com/bzelaznicki/bzCommerce/actions/workflows/frontend-ci.yml/badge.svg)](https://github.com/bzelaznicki/bzCommerce/actions/workflows/frontend-ci.yml)

bzCommerce is a modern e-commerce platform featuring a robust Go backend and a Next.js frontend. It provides a seamless shopping experience for users and powerful management tools for administrators. The platform is designed for speed, security, and scalability.

## Features

- **User-Friendly Interface**: Intuitive design for customers to browse and purchase products.
- **Admin Dashboard**: Manage categories, products, and users with ease.
- **Authentication**: Secure user registration and login system.
- **Database Integration**: Efficient data handling using SQL (PostgreSQL recommended).
- **Docker Support**: Simplified deployment with Docker and Docker Compose.

## Tech Stack

- **Frontend**: Next.js (React), TypeScript, HTML, CSS
- **Backend**: Go
- **Database**: PostgreSQL
- **Containerization**: Docker

---

## Monorepo Structure

- `/frontend` – Next.js frontend application
- `/backend` – Go backend API and server

---

## Setup Instructions

### 1. Clone the Repository
```bash
git clone https://github.com/bzelaznicki/bzCommerce.git
cd bzCommerce
```

### 2. Backend Setup
- Ensure you have Go 1.23+ installed.
- Install dependencies:
  ```bash
  cd backend
  go mod tidy
  ```
- Set up the database:
  - Create a PostgreSQL database.
  - Update the database connection details in `sqlc.yaml`.
  - Run migrations:
    ```bash
    ../scripts/migrateup.sh
    ```
- Set up environment variables:
  - Copy `.env.example` to `.env` and update values as needed.
- Run the backend server:
  ```bash
  go run main.go
  ```

### 3. Frontend Setup
- Ensure you have Node.js 18+ and npm installed.
- Install dependencies:
  ```bash
  cd ../frontend
  npm install
  ```
- Start the development server:
  ```bash
  npm run dev
  ```
- The frontend will be available at `http://localhost:3000` by default.

### 4. Using Docker (Optional)
- Build and start the entire stack using Docker Compose:
  ```bash
  docker-compose up --build
  ```

---

## Usage

- Access the shop at `http://localhost:3000` (frontend).
- Backend API runs at `http://localhost:8080` (default, configurable in `.env`).
- Admin dashboard is available at `/admin` (credentials required).

---

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes and push the branch.
4. Open a pull request.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

