#!/bin/sh


export VAULT_ADDR=http://localhost:8200
export VAULT_TOKEN=default-token
vault login
myroot

# vault 1.1
vault kv put secret/prod user=richard
vault kv get secret/prod/user
vault kv get -format=json secret/prod

vault kv put secret/prod/api.key value=123456789
vault kv put secret/prod/api.token value=abcdefghijklmnopqrstuvwxyz

# vault 0.7.3
vault list secret/prod
vault write secret/prod user=richard
vault read secret/prod
