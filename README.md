# vault-env
An application to log into vault and render a secrets file for other applications to use.

## What is VaultEnv

Vaultenv is an application that will login into vault, read a set of secrets from vault and render them before applications start. 


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