# Order Service

Order Service is a backend service for managing orders, utilizing PostgreSQL for database storage, in-memory caching, and NATS Streaming for message processing.

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [Running Tests](#running-tests)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## Installation

1. **Clone the repository:**

    ```sh
    git clone https://github.com/F-1X/WB-L0.git
    cd order-service
    ```

2. **Install dependencies:**

    Ensure you have Go installed, then run:

    ```sh
    go mod tidy
    ```

3. **Set up PostgreSQL and NATS Streaming:**

    Ensure you have PostgreSQL installed and running. Create a database for the project.
    ```sh
    make docker-up
    make mig-up
    ```

    

4. **Run an application:**

     ```sh
    make run-server
    ```

## Configuration

Set the following environment variables in a `.env` file or export them directly in your terminal:

```env
DB_USER=user
DB_HOST=localhost
DB_PORT=5432
DB_PASSWORD=example
DB_NAME=orderdb
MIGRATIONS_PATH=backend/internal/database/migrate/migrations
DATABASE_URL=postgres://user:example@localhost:5432/orderdb
NATS_URL=nats://nats:4222
CONFIG_PATH=./backend/internal/config/config.yml
FRONTEND_PATH=./frontend/static
TEST_DB_URI=postgres://test:test@localhost:5432/orderdb
```
## Databse structure of tables
**orders**

| Column              | Type           | Constraints         |
|---------------------|----------------|---------------------|
| order_uid           | VARCHAR(50)    | PRIMARY KEY, UNIQUE |
| track_number        | VARCHAR(50)    | UNIQUE              |
| entry               | VARCHAR(50)    |                     |
| locale              | VARCHAR(10)    |                     |
| internal_signature  | VARCHAR(50)    |                     |
| customer_id         | VARCHAR(50)    |                     |
| delivery_service    | VARCHAR(50)    |                     |
| shardkey            | VARCHAR(10)    |                     |
| sm_id               | BIGINT         |                     |
| date_created        | TIMESTAMP      |                     |
| oof_shard           | VARCHAR(10)    |                     |

**delivery**

| Column      | Type           | Constraints                                |
|-------------|----------------|--------------------------------------------|
| id          | SERIAL         | PRIMARY KEY                                |
| order_uid   | VARCHAR(50)    | FOREIGN KEY REFERENCES orders(order_uid)   |
| name        | VARCHAR(100)   |                                            |
| phone       | VARCHAR(20)    |                                            |
| zip         | VARCHAR(20)    |                                            |
| city        | VARCHAR(100)   |                                            |
| address     | VARCHAR(200)   |                                            |
| region      | VARCHAR(100)   |                                            |
| email       | VARCHAR(100)   |                                            |

**payment**

| Column        | Type           | Constraints                                |
|---------------|----------------|--------------------------------------------|
| id            | SERIAL         | PRIMARY KEY                                |
| transaction   | VARCHAR(50)    | FOREIGN KEY REFERENCES orders(order_uid)   |
| request_id    | VARCHAR(100)   |                                            |
| currency      | VARCHAR(10)    |                                            |
| provider      | VARCHAR(100)   |                                            |
| amount        | INT            |                                            |
| payment_dt    | INT            |                                            |
| bank          | VARCHAR(100)   |                                            |
| delivery_cost | INT            |                                            |
| goods_total   | INT            |                                            |
| custom_fee    | INT            |                                            |

**items**

| Column       | Type           | Constraints                                    |
|--------------|----------------|------------------------------------------------|
| id           | SERIAL         | PRIMARY KEY                                    |
| chrt_id      | INT            |                                                |
| track_number | VARCHAR(50)    | FOREIGN KEY REFERENCES orders(track_number)    |
| price        | INT            |                                                |
| rid          | VARCHAR(100)   |                                                |
| name         | VARCHAR(200)   |                                                |
| sale         | INT            |                                                |
| size         | VARCHAR(20)    |                                                |
| total_price  | INT            |                                                |
| nm_id        | INT            |                                                |
| brand        | VARCHAR(100)   |                                                |
| status       | INT            |                                                |
