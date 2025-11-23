# Contributing to Community Issue Mapper

First off, thank you for considering contributing to Community Issue Mapper! It's people like you that make this project a great tool for communities.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When you create a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Provide specific examples** to demonstrate the steps
- **Describe the behavior you observed** and what you expected
- **Include screenshots** if relevant
- **Include your environment details** (OS, browser, Go version, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide a detailed description of the suggested enhancement**
- **Explain why this enhancement would be useful**
- **List some examples** of how it would be used

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** with clear, descriptive commits
3. **Add tests** if applicable
4. **Ensure all tests pass**: `make test`
5. **Update documentation** if needed
6. **Submit a pull request**

## Development Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Docker and Docker Compose (optional)
- Make

### Local Development

1. **Clone your fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/community-issue-mapper.git
   cd community-issue-mapper
   ```

2. **Install dependencies**
   ```bash
   cd backend
   make deps
   ```

3. **Set up database**
   ```bash
   createdb community_issues
   make migrate-up
   ```

4. **Run the application**
   ```bash
   make run
   ```

5. **Run tests**
   ```bash
   make test
   ```

### Development with Docker

```bash
docker-compose up -d
```

## Coding Standards

### Go Code

- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format your code
- Run `golangci-lint run` before committing
- Write meaningful variable and function names
- Add comments for exported functions and types
- Keep functions small and focused

### JavaScript Code

- Use ES6+ features
- Use meaningful variable names
- Add JSDoc comments for functions
- Keep functions pure when possible
- Avoid global state

### Git Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters
- Reference issues and pull requests when relevant

Example:
```
Add image optimization for uploads

- Resize images to max 1920px width
- Convert to WebP format when supported
- Add progressive JPEG encoding
- Update upload handler with optimization

Closes #123
```

## Project Structure

```
community-issue-mapper/
├── backend/                   # Go backend
│   ├── cmd/                  # Application entry points
│   ├── internal/             # Private application code
│   │   ├── api/             # API handlers and routing
│   │   ├── config/          # Configuration
│   │   ├── database/        # Database connection
│   │   ├── models/          # Data models
│   │   ├── repository/      # Data access layer
│   │   └── service/         # Business logic
│   └── migrations/          # Database migrations
├── frontend/                 # Frontend code
│   ├── css/                 # Stylesheets
│   └── js/                  # JavaScript
└── docs/                    # Documentation
```

## Testing Guidelines

### Writing Tests

- Write tests for new features
- Update tests when modifying existing features
- Aim for good test coverage (>70%)
- Use table-driven tests when appropriate
- Test edge cases and error conditions

Example Go test:
```go
func TestCreateIssue(t *testing.T) {
    tests := []struct {
        name    string
        input   *models.CreateIssueRequest
        wantErr bool
    }{
        {
            name: "valid issue",
            input: &models.CreateIssueRequest{
                Title:       "Test Issue",
                Description: "Test Description",
                Category:    models.CategoryPothole,
                Latitude:    37.7749,
                Longitude:   -122.4194,
            },
            wantErr: false,
        },
        {
            name: "invalid latitude",
            input: &models.CreateIssueRequest{
                Title:       "Test Issue",
                Description: "Test Description",
                Category:    models.CategoryPothole,
                Latitude:    91.0,
                Longitude:   -122.4194,
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateIssue(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateIssue() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## Documentation

- Update README.md if you change functionality
- Update API.md for API changes
- Add inline code comments for complex logic
- Update CHANGELOG.md for notable changes

## Review Process

1. At least one maintainer must approve the PR
2. All CI checks must pass
3. Code must follow project conventions
4. Documentation must be updated if needed

## Community

- Join discussions in GitHub Issues
- Help others with their questions
- Share your use cases and ideas
- Be respectful and constructive

## Recognition

Contributors will be recognized in:
- GitHub contributors page
- CHANGELOG.md for significant contributions
- Special mentions for major features

## Questions?

Feel free to open an issue with the `question` label or reach out to the maintainers.

Thank you for contributing to Community Issue Mapper!
