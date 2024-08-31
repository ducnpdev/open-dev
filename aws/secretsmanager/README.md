
# aws secrets manager

aws secretsmanager list-secrets

aws secretsmanager list-secrets \
    --query "SecretList[?Tags[?Key=='aws secretsmanager list-secrets --filter Key="asm-apse1-cicdhdbank/dev/nhdt/fakeapp/PvoilSecret"']]"

aws secretsmanager list-secrets --filter Key="asm-apse1-cicdhdbank/dev/nhdt/fakeapp/PvoilSecret"

aws secretsmanager create-secret --name asm-apse1-cicdhdbank/dev/nhdt/fakeapp/DatabasePassword \
                                --description "My test secret created with the CLI." \
                                --secret-string "{\"user\":\"diegor\",\"password\":\"EXAMPLE-PASSWORD\"}"
                        


aws secretsmanager create-secret --name asm-apse1-cicdhdbank/dev/nhdt/fakeapp/PvoilSecret \
                                --description "My test secret created with the CLI." \
                                --secret-string "{\"user\":\"diegor\",\"password\":\"EXAMPLE-PASSWORD\"}"

aws secretsmanager list-secrets | jq '.SecretList[] | select(.Name | contains("fakeapp"))'
