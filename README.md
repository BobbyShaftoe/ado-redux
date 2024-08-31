
# HTTP Server Project

This project is an HTTP server built using Go. It interacts with Azure DevOps to fetch and display project and repository data. The server also validates users and serves static files.

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Endpoints](#endpoints)
- [Logging](#logging)
- [Contributing](#contributing)
- [License](#license)

## Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/yourusername/HTTP_Server.git
   cd HTTP_Server
   ```

2. **Install dependencies:**
   Ensure you have Go installed. Then, run:
   ```sh
   go mod tidy
   ```

3. **Set up the database:**
   Ensure you have PostgreSQL installed and running. Create a database and user as specified in the `main.go` file.

## Configuration

Create a `.env` file in the root directory with the following environment variables:
```
USER=your_azure_devops_user
PROJECT=your_azure_devops_project
```

## Usage

Run the server using:
```sh
go run main.go
```

The server will start at `http://localhost:8080`.

## Endpoints

- **GET /**: Renders the index page.
- **GET /dashboard**: Renders the dashboard page.
- **GET /repositories**: Renders the repositories page.
- **POST /dashboard-update**: Updates the dashboard with a new project.
- **POST /repositories-update**: Updates the repositories with a new project.
- **GET /search**: Handles search functionality.

## Logging

The project uses `slog` for logging. Logs are output in JSON format to `stdout`.

## Contributing

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Create a new Pull Request.

## License

This project is licensed under the MIT License.

