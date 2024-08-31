# Project Name

## Summary

This project is a web application that integrates with Azure DevOps to aggregate and report statuses for pipelines,
repositories, and commits. It includes functionally for running builds, finding work item relations and searching commits, repos etc.

## Git Branches
All features are developed in feature branches and merged into the develop branch.
Periodically, the develop branch is merged into the main branch.

## Project Requirements

- Go 1.16 or later
- Templ and Tailwind CSS tools installed
- Azure DevOps account and PAT (Personal Access Token)
- Environment variables for configuration

## Programming Languages and Components

- Go + Templ
- HTMX
- Alpine.js
- Basic JavaScript
- Tailwind CSS


## Configuration

Create an `.env` file in the root directory or export the following environment variables:
```
ADO_ORG="Name Of ADO Organization"
ADO_DEFAULT_PROJECT=YourProjectName
ADO_DEFAULT_USER=your.ado.email@alt.net
AZURE_TOKEN=*********************************
```

## Usage

First examine the Makefile to see the available commands and how to develop and build the project.

**Run the server using:**
```sh
make run
```
OR
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



## Project Architecture

The project is structured into several packages and files:

- `main.go`: Entry point of the application.
- `model`: Contains data models and business logic.
- `helpers`: Contains helper functions and configurations.
- `static`: Contains static files like CSS and JavaScript.

### Packages and Files

#### `main.go`

- **Description**: Entry point of the application.

#### `model/global.go`

- **GlobalState**: Manages the global state of the application.
   - `User`: Current user.
   - `AdditionalUsers`: List of additional users.
   - `UserValidated`: Boolean indicating if the user is validated.
   - `Projects`: List of projects.
   - `CurrentProject`: Current project.

- **Functions**:
   - `NewGlobalState(user string, projects []string) *GlobalState`: Initializes a new global state.
   - `UpdateGlobalState(user string, projects []string)`: Updates the global state.
   - `UpdateGlobalStateProjects(projects []string)`: Updates the projects in the global state.
   - `UpdateGlobalStateProject(project string)`: Updates the current project in the global state.
   - `UpdateGlobalStateUser(user string)`: Updates the user in the global state.
   - `UpdateUserValidated(status bool)`: Updates the user validation status.

#### `model/repositories.go`

- **RepositoriesData**: Holds data for projects and repositories.
   - `Projects`: List of projects.
   - `Repos`: List of repositories.

#### `model/dashboard.go`

- **HelloData**: Holds a name for greeting.
   - `Name`: Name for greeting.

- **DashboardData**: Holds data for the dashboard.
   - `Projects`: List of projects.
   - `Repos`: List of repositories.
   - `Commits`: List of commits.

#### `model/ado.go`

- **ADOConnectionInfo**: Holds Azure DevOps connection information.
   - `ConnectionUrl`: URL for the connection.
   - `ConnectionPAT`: Personal Access Token for the connection.

- **GitCommitsCriteria**: Criteria for fetching Git commits.
   - `RepositoryId`: ID of the repository.
   - `Author`: Author of the commits.
   - `User`: User who made the commits.
   - `FromDate`: Start date for the commits.
   - `Version`: Version of the commits.
   - `VersionType`: Type of the version.
   - `Skip`: Number of commits to skip.
   - `Top`: Number of commits to fetch.

- **GitCommitItem**: Holds information about a Git commit.
   - `Repository`: Name of the repository.
   - `CommitInfo`: List of commit information.

#### `helpers/config/config.go`

- **EnvVars**: Holds environment variables.
   - `DBPass`: Database password.
   - `PAT`: Personal Access Token.
   - `ORGANIZATION`: Azure DevOps organization.
   - `PROJECT`: Default project.
   - `USER`: Default user.

- **Functions**:
   - `New() *EnvVars`: Initializes a new set of environment variables.
   - `getEnv()`: Fetches environment variables.

#### `helpers/ado/ado.go`

- **GitRepo**: Holds information about a Git repository.
   - `Name`: Name of the repository.
   - `Id`: ID of the repository.
   - `Url`: URL of the repository.
   - `WebUrl`: Web URL of the repository.
   - `Links`: Links related to the repository.
   - `DefaultBranch`: Default branch of the repository.

- **Functions**:
   - `ReturnProjects(responseValue *core.GetProjectsResponseValue) []string`: Returns a list of project names.
   - `ReturnGitRepos(responseValue *[]git.GitRepository) []GitRepo`: Returns a list of Git repositories.
   - `ReturnGitRepoNames(gitRepositories *[]git.GitRepository) []string`: Returns a list of Git repository names.
   - `ValidateUser(user string, userGraph *[]graph.GraphUser) bool`: Validates a user.

#### `static/light-dark-mode.js`

- **Description**: Manages the light/dark mode toggle.
- **Functions**:
   - Toggles the light/dark mode based on user preference stored in local storage.

## Logging

The project uses `slog` for logging. Logs are output in JSON format to `stdout`.


### JavaScript Components

- **Alpine.js**: Used for widgets/page components.
- **Light/Dark Mode Toggle**: Manages the light/dark mode based on user preference stored in local storage.


## Contributing

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -am 'Add new feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Create a new Pull Request.

## License

This project is licensed under the MIT License.