# bzCommerce

bzCommerce is an e-commerce platform designed to provide a seamless shopping experience for users and robust management tools for administrators. Built with Go, it leverages modern web technologies to deliver a fast, secure, and scalable solution for online stores.

## Features

- **User-Friendly Interface**: Intuitive design for customers to browse and purchase products.
- **Admin Dashboard**: Manage categories, products, and users with ease.
- **Authentication**: Secure user registration and login system.
- **Database Integration**: Efficient data handling using SQL.
- **Docker Support**: Simplified deployment with Docker and Docker Compose.

## Tech Stack

- **Backend**: Go
- **Database**: SQL (PostgreSQL recommended)
- **Frontend**: HTML, CSS
- **Containerization**: Docker

## Setup Instructions

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/bzelaznicki/bzCommerce.git
   cd bzCommerce
   ```

2. **Install Dependencies**:
   Ensure you have Go 1.23 installed, then run:
   ```bash
   go mod tidy
   ```

3. **Set Up the Database**:
   - Create a PostgreSQL database.
   - Update the database connection details in `sqlc.yaml`.
   - Run migrations:
     ```bash
     ./scripts/migrateup.sh
     ```

4. **Set Up Environment Variables**:
   - Copy the `.env.example` file to `.env`:
     ```bash
     cp .env.example .env
     ```
   - Update the `.env` file with your specific configuration values, including setting the desired port for the application.

5. **Run the Application**:
   ```bash
   go run main.go
   ```

6. **Using Docker (Optional)**:
   - Build and start the application using Docker Compose:
     ```bash
     docker-compose up --build
     ```

## Usage

- Access the homepage at `http://localhost:8080`.
- Admin dashboard is available at `http://localhost:8080/admin` (credentials required).
- Ensure the port specified in the `.env` file matches the one you use to access the application (e.g., `http://localhost:8080`).

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes and push the branch.
4. Open a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.