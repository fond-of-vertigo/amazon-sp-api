# Scripts

## Download and extract Amazon SP API models

Using the `pull_and_extract_model` script, you are able to extract collected 
models from a openapi.{yaml|json} or url to a single file.

Example:
```sh
./pull_and_extract_model.sh https://raw.githubusercontent.com/amzn/selling-partner-api-models/main/models/finances-api-model/financesV0.json model.go
```
