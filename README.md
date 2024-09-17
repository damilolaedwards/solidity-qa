# Audit Assistant

**Audit Assistant** is a powerful, lightweight command-line tool designed to streamline the process of conducting and managing security audits. It helps security professionals and developers identify vulnerabilities, manage audit findings, and generate comprehensive audit reports. With customizable templates and the ability to initialize new audit projects, this tool is ideal for conducting thorough audits on both local and on-chain contracts.

## Features

- **Project Initialization:** Quickly generate and set up audit project configurations with customizable parameters.
- **Audit Report Generation:** Automatically generate reports based on parsed contract data for different audit types.
- **Manage Findings:** Track vulnerabilities and manage audit recommendations.
- **On-Chain Support:** Supports auditing on-chain contracts via Etherscan.

## Prerequisites

Before downloading medusa, you will need to have `crytic-compile` and `slither` installed.

- Installation instructions for `crytic-compile` can be found [here](https://github.com/crytic/crytic-compile).
- Installation instructions for `slither` can be found [here](https://github.com/crytic/slither).

`crytic-compile` and `slither` require a Python environment. Installation instructions for Python can be found [here](https://www.python.org/downloads/).

## Installation

1. Clone the repository from the source:

   ```bash
   git clone <repository-url>
   ```

2. Navigate into the project directory:

   ```bash
   cd assistant
   ```

3. Build the binary:

   ```bash
   go build -o assistant
   ```

4. Add the binary to your system's PATH to access it globally:

   - On Linux or macOS, add the following line to your `~/.bashrc`, `~/.bash_profile`, or `~/.zshrc`:
     ```bash
     export PATH=$PATH:/path/to/assistant
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
assistant init
```

This will generate a configuration file called `assistant.json` in the current directory.

#### With Options

```bash
assistant init --out="config.json" --name="my-audit" --port="9000" --target-contracts-dir="./contracts" --test-contracts-dir="./tests"
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
assistant start
```

#### With Target Contracts Directory

```bash
assistant start ./path/to/contracts
```

#### With Options

```bash
assistant start --config="config.json" --onchain --address="0xABC123" --api-key="$ETHERSCAN_API_KEY"
```

This will use the specified configuration file and spin up the audit process, fetching contract details from Etherscan if the `onchain` flag is set.

**Flags:**

- `--config`: Path to the configuration file.
- `--onchain`: Specifies if the contract is an on-chain contract rather than a local project.
- `--address`: Address of the on-chain contract to be analyzed.
- `--network`: Network of the on-chain contract to be analyzed.
  - Supported networks: `mainnet`, `arbitrum`, `optimism`, `polygon`, `bsc`, `avalanche`, `fantom`
  - Default: `mainnet`
- `--api-key`: API key for fetching on-chain contract data.

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
  "includeInterfaces": false, // Whether interfaces will be included in the slither output
  "includeAbstract": false, // Whether abstract contracts will be included in the slither output
  "includeLibraries": false // Whether libraries will be included in the slither output
}
```

## License

Audit Assistant is released under the MIT License.
