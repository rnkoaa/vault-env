# vault-env
An a CLI application that can render secrets file for other applications to use.

## Capabilities

Vaultenv is an application that will login into vault, read a set of secrets from vault and render them before applications start. 

Installation

From source:

## Installation

### Requirements

Go `1.12.6` or higher which can be downloaded from [here](https://golang.org/dl/).

### From source:

1.  Run `git clone <vault-env repo>`
2.  Run `GO111MODULE=on go mod vendor`
3.  Run `cd cmd/vaultenv` 
3.  Run `go build -o vault-env .`

### From Artifactory:

Coming soon.

## Methods of Configuration

### CLI

```
./vault-env --help  
As applications are deployed to in containers, they will need their secrets. 
In order to prevent secret sprawl and prevent secret leakage, 
this application takes authentication configs from the orchestration environment, 
logs into vault and generates the secret properties before the application starts up.

Usage:
  vault-env [command]

Available Commands:
  help        Help about any command
  render      render values from vault

Flags:
      --config string   config file (default is config.yml)
  -h, --help            help for vault-env

Use "vault-env [command] --help" for more information about a command.
```

## Configuration

The application can be configured in several ways. 

1. using a default `runtime.yml` file 
2. a secondary file called `config.yml` which can override values in the `runtime.yml` file.
3. Also  through command line flags `--config config` which will override any values previous set.
4. It can also read environment variables which can be overriden in any of the above configs.


### Example template file

Below is a sample configuration file with will be used to render a set of secrets.

```
api:
  key: api.key # key from vault this will be rendered along with the keyPrefix
  token: api.token # key from vault this will be rendered along with the keyPrefix
cassandra:
  username: cassandra.username
  password: cassandra.password
  user: secret/github/prod/cassandra.password # this will use the entire key since it starts with 'secret/'
```

## Usage

```sh
# This reads runtime.yml file and overrides any config from config.yaml file
$ ./vault-env render --config config.yaml 

# This reads runtime.yml file without any external overrides. 
$ ./vault-env render
```

## Authentication

Currently, only two vault authentication modes are supported.

- `token` 
- `LDAP` 

## Known Issues

## Contributing

Contributions in the form of patches, PR's and Issues are welcome. All submissions, including submissions by project members, require review. We use GitHub pull requests for this purpose. 

## Developing

`vault-env` was built so applications can easily retrieve their secrets from a central location without baking it into the applications.

1. Once built, create a `runtime.yml` file. 
2. Add necessary configs into the `runtime.yml` file to point to the correct `vault_team` and also vault `keyPrefix` template. Please note that the key prefix is can either be of the following formats: 

    - raw: `github/dev` 
    - go templates: `{{.VaultTeam}}/{{CloudEnvironment}}`. 
    - mixed: `github/{{CloudEnvironment}}`
    
    However, if there is no keyPrefix defined, the key generated should include the full path in vault such as `secret/github/dev/{{key}}` for `Vault KV version 1` and `secret/data/github/dev/{{dev}}` for `Vault KV version 2`.