# Solidity-QA

**Solidity-QA** is a lightweight conversational AI code review tool designed to streamline the process of conducting and managing security audits. It leverages on the use AI to scan and identify code smells and vulnerabilities in codebases as well as generate comprehensive issue reports. With customizable configuration, it supports the use of both supports chatGPT and Claude APIs

## Features

- **Web based interface:** Intuitive chat user interface with color coded elements to enhance readability.
- **Context Switch:** Seamlessly switch context between chatGPT and Claude within the same conversation flow.
- **Report Generation:** Generate issue report based of predefined template.
- **On-Chain Support:** Supports auditing on-chain contracts via Etherscan.

## Prerequisites

Before downloading solidity-qa, you will need to have `crytic-compile` and `slither` installed.

- Installation instructions for `crytic-compile` can be found [here](https://github.com/crytic/crytic-compile).
- Installation instructions for `slither` can be found [here](https://github.com/crytic/slither).

`crytic-compile` and `slither` require a Python environment. Installation instructions for Python can be found [here](https://www.python.org/downloads/).

- Note: Python version must be `3.10` or greater

## Installation

1. Clone the repository from the source:

   ```bash
   git clone https://github.com/damilolaedwards/solidity-qa
   ```

2. Navigate into the project directory:

   ```bash
   cd solidity-qa
   ```

3. Build the binary:

   ```bash
   go build -o solidity-qa
   ```

4. Add the binary to your system's PATH to access it globally:

   - On Linux or macOS, add the following line to your `~/.bashrc`, `~/.bash_profile`, or `~/.zshrc`:
     ```bash
     export PATH=$PATH:/path/to/solidity-qa
     ```
   - On Windows, add the binary directory to your system's environment variables.

5. Add your API keys. You can use the sample file as a template:

   ```bash
   cp api_keys.sample.sh api_keys.sh
   ```

   Update the keys in the file, then source the file to add your keys to your environment:

   ```bash
   source api_keys.sh
   ```

## Usage

### 1. Initialize config

The `init` command generates a configuration file for a new audit project. You can specify various options such as project name, port, and contract directories.

#### Basic Usage

```bash
solidity-qa init
```

This will generate a configuration file called `assistant.json` in the current directory.

#### With Options

```bash
solidity-qa init --out="config.json" --name="my-audit" --port="9000" --target-contracts-dir="./contracts" --test-contracts-dir="./tests"
```

This will generate a `config.json` file with custom project details.

**Flags:**

- `--out`: Output path for the config file (default is `assistant.json`).
- `--name`: Name of the project.
- `--port`: Port number for the API (default is `8080`).
- `--target`: Directory containing the contracts to be audited.
- `--test-dir`: Directory containing the test contracts.

### 2. Start the Audit

Once the project has been initialized and the configuration file is populated, you can start the audit process.

#### Basic Usage

```bash
solidity-qa start
```

#### With Target Contracts Directory

```bash
solidity-qa start ./path/to/contracts
```

#### With Options

```bash
solidity-qa start --config="config.json" --onchain --exclude-interfaces --address="0xABC123" --api-key="$ETHERSCAN_API_KEY"
```

This will use the specified configuration file and spin up the session, fetching contract source code from Etherscan if the `onchain` flag is set.

**Flags:**

- `--config`: Path to the configuration file.
- `--port`: The port the local server will be served on.
- `--host`: Whether the local server will be hosted using ngrok (requires the NGROK_AUTHTOKEN env variable).
- `--slither-args`: Extra arguments to be passed to slither.
- `--onchain`: Specifies if the contract is an on-chain contract rather than a local project.
- `--address`: Address of the on-chain contract to be analyzed.
- `--network`: Network of the on-chain contract to be analyzed.
  - Supported networks: `mainnet`, `arbitrum`, `optimism`, `polygon`, `bsc`, `avalanche`, `fantom`
  - Default: `mainnet`
- `--api-key`: API key for fetching on-chain contract data.
- `--exclude-interfaces`: Specifies if interfaces should be excluded from the analysis.

### Example Config File

A sample `assistant.json` config file looks like this:

```json
{
  "name": "example", // The project name
  "targetContracts": {
    "directory": "", // The directory path relative to the project root
    "excludePaths": [] // Paths that should be excluded when parsing the directory
  },
  "testContracts": {
    "directory": "", // The directory path relative to the project root
    "excludePaths": [] // Paths that should be excluded when parsing the directory
  }, // The directory that holds the test contracts
  "port": 8080, // The port that the API will be running on
  "host": false, // Whether the local server should be hosted using ngrok
  "includeInterfaces": false, // Whether interfaces will be included in the slither output
  "includeAbstract": false, // Whether abstract contracts will be included in the slither output
  "includeLibraries": false, // Whether libraries will be included in the slither output,
  "slitherArgs": {} // Extra arguments to be passed to slither
}
```

#### Passing extra slither args

In cases where you need to pass extra arguments to slither, the `slitherArgs` config option can be used to provide the extra arguments. Examples:

- As command-line argument:

  ```bash
  solidity-qa start --slither-args='{ "compile_force_framework": "hardhat" }'
  ```

- In config:

  ```json
  {
    ... // Other config parameters
    "slitherArgs": {
      "compile_force_framework": "hardhat",
      ...
    }
  }
  ```

## License

Solidity-QA is released under the MIT License.
