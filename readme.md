# GitHubAuto

GitHubAuto is a command-line tool written in Go that interacts with the GitHub API to perform various operations such as searching repositories, querying issues, checking notifications, and more.

## Features

- Save GitHub Personal Access Token
- Fetch all repositories for the authenticated user
- Search for a GitHub repository
- Show details of a specific GitHub repository
- Query issues for a specific repository
- Check notifications for the user

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/itsfuad/githubauto.git
    ```
2. Navigate to the project directory:
    ```sh
    cd githubauto
    ```
3. Build the project:
    ```sh
    go build -o githubauto
    ```

## Usage

Run the executable with the desired options:

```sh
./githubauto [options]
```

### Options

- `-save-token` : Save GitHub Personal Access Token
- `-all-repo` : Fetch all repositories for the authenticated user
- `-search-repo <repo>` : Search for a GitHub repository
- `-show-repo <repo>` : Show details of a GitHub repository
- `-query-issues <repo>` : Query issues for a specific repository
- `-notify` : Check notifications for the user

### Examples

Save GitHub Token:
```sh
./githubauto -save-token
```

Fetch all repositories:
```sh
./githubauto -all-repo
```

Search for a repository:
```sh
./githubauto -search-repo go
```

Show repository details:
```sh
./githubauto -show-repo yourusername/repo
```

Query issues in a repository:
```sh
./githubauto -query-issues yourusername/repo
```

Check notifications:
```sh
./githubauto -notify
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## Acknowledgements

- [Go](https://golang.org/)
- [GitHub API](https://docs.github.com/en/rest)
