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
 export VAULT_TOKEN=default-token
 export VAULT_ADDR=http://localhost:8200
 vault write secret/github/prod/cassandra.username value=cassandra_manager
 vault write secret/github/prod/cassandra.password value=cassandra_password
 vault write secret/github/prod/api.token value=1234567890
 vault write secret/github/prod/api.key value=abcdefghijklmnopqrstuvwxyz
