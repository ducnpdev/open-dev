
# aws code artifact
## aws-artifact push package
- create domain
aws codeartifact create-domain --domain [domain-name]
- create repo
aws codeartifact create-repository --domain [domain-name] --domain-owner [account-id] --repository [repo-name]
- login
aws codeartifact login --tool npm --repository [repo-name] --domain [domain-name] --domain-owner [account-id]
- get list package
aws codeartifact list-packages --domain [domain-name] --repository [repo-name]
- push package
aws codeartifact publish-package-version  --domain [domain-name] --repository [repo-name] --format [format-type] --namespace [name-ns] --asset-sha256 [hash] --package [package] --package-version [version] --asset-content [path-file] --asset-name [asset]
- delete repo
aws codeartifact delete-repository --domain [domain-name] --domain-owner [account-id] --repository [repo-name]
- delete domain
aws codeartifact delete-domain --domain [domain-name] --domain-owner [account-id]

  - explain
    - --asset-content: đường dẫn file cần push lên
    - --asset-sha256: hash sha246 của file
        `openssl dgst -sha256 abc.txt`
    - --asset-name: tên của artifact khi cần pull về
    - --format: type
    - --namespace: namespace

### Example:
- create domain
aws codeartifact create-domain --domain open-dev-domain

- create repo
aws codeartifact create-repository --domain open-dev-domain --domain-owner 06***48 --repository open-dev-repo

- login
aws codeartifact login --tool npm --repository open-dev-repo --domain open-dev-domain --domain-owner 06***48

- get list package
aws codeartifact list-packages --domain open-dev-domain --repository open-dev-repo

- push package
aws codeartifact publish-package-version  --domain open-dev-domain --repository open-dev-repo --format generic --namespace abc --asset-sha256         e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 --package abc --package-version 1.0.0 --asset-content abc.txt --asset-name abc 

- delete repo
aws codeartifact delete-repository --domain open-dev-domain --domain-owner 06***48 --repository open-dev-repo

- delete domain
aws codeartifact delete-domain --domain open-dev-domain --domain-owner 06***48
