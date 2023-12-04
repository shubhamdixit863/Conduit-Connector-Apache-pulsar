
# Apache Pulsar Conduit Connector

This project provides a Conduit connector for Apache Pulsar, enabling seamless integration with Conduit's data streaming capabilities. Developed in Go, this connector is designed to work efficiently with Apache Pulsar, a distributed messaging and streaming platform.

## Features

- **Pulsar Client Integration**: Fully integrated with the Apache Pulsar client, ensuring reliable and efficient messaging.
- **Source and Destination Support**: Skeleton code for both source (reading from Pulsar) and destination (writing to Pulsar) are provided.
- **Configurable**: Easy to configure for different environments and Pulsar clusters.
- **Unit Tests**: Example unit tests for validating the functionality of the connector.
- **CI/CD Integration**: GitHub Actions workflows for automated building, testing, and linting.

## Getting Started

### Prerequisites

- Go (version 1.21)
- Access to an Apache Pulsar instance

### Usage

1. **Use this Template**: Click "Use this template" on the repository's main page.
2. **Clone the Repository**: After setting up your repository, clone it to your local machine.
3. **Initial Setup**: Run `./setup.sh <module name>` to initialize the module (e.g., `./setup.sh github.com/your-org/conduit-connector-pulsar`).
4. **Configuration**: Modify the `config.go` file to set up general and Pulsar-specific configurations.

### Building and Running

- To build the connector, use `make build`.
- Run unit tests with `make test`.

## Repository Setup Recommendations

- **Pull Requests**: Use squash merging, enable auto-merge, and set up automatic branch deletion post-merge.
- **Branch Protection**: Enforce pull requests and approvals, require status checks, and conversation resolution before merging.
- **Actions**: Enable all actions and allow Actions to manage pull requests, especially for Dependabot configurations.

## Connector Specification

- The `spec.go` file includes the configuration schema for the Pulsar connector.

## Contributing

Contributions are welcome! Please follow the guidelines in the `CONTRIBUTING.md` file.

## License

This project is licensed under the [MIT License](LICENSE).


