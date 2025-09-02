# Contributing to OtterScale

Thank you for your interest in contributing to OtterScale! We welcome contributions from the community and appreciate your help in making this project better.

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally
3. Create a new branch for your feature or bug fix
4. Make your changes
5. Test your changes thoroughly
6. Submit a pull request

## Development Setup

### Prerequisites

Before starting development, ensure you have the following installed:

- [Node.js](https://nodejs.org/) (v18 or higher)
- [pnpm](https://pnpm.io/) package manager
- [Go](https://golang.org/) (v1.21 or higher)
- [Git](https://git-scm.com/)

### Setup Instructions

1. **Clone your fork**

   ```bash
   git clone https://github.com/your-username/otterscale.git
   cd otterscale
   ```

2. **Install dependencies**

   ```bash
   # Frontend dependencies
   pnpm install

   # Backend dependencies
   go mod download
   ```

3. **Run tests**

   ```bash
   # Frontend tests
   pnpm test

   # Backend tests
   go test ./...
   ```

4. **Start development servers**

   ```bash
   # Frontend development server
   pnpm run dev

   # Backend development server (in a separate terminal)
   go run main.go
   ```

### Verification

After setup, verify everything works:

- Frontend should be accessible at `http://localhost:3000`
- Backend API should be running on `http://localhost:8299`
- All tests should pass without errors

## Code Style

- Follow existing code formatting and style
- Use meaningful variable and function names
- Add comments for complex logic
- Ensure your code passes linting checks

## Commit Messages

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification. Please format your commit messages as:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code formatting changes
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples

```
feat(auth): add OAuth2 login support
fix(ui): correct button alignment in modal
docs: update installation instructions
```

## Pull Request Process

1. Update documentation if needed
2. Add tests for new functionality
3. Ensure all tests pass
4. Update the README.md if necessary
5. Request review from maintainers

## Reporting Issues

When reporting issues, please include:

- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Node.js version, etc.)
- Screenshots if applicable

## Code of Conduct

Please be respectful and inclusive in all interactions. We aim to create a welcoming environment for all contributors.

## Questions?

If you have questions about contributing, feel free to open an issue or reach out to the maintainers.

Thank you for contributing to OtterScale!
