#!/bin/sh

vault 1.1

export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=default-token
vault login
myroot

vault kv put secret/prod user=richard
vault kv get secret/prod/user
vault kv get -format=json secret/prod


vault 0.7.3
vault list secret/prod
vault write secret/prod user=richard
vault read secret/prod