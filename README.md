# add-email
OpenFaaS Go function to add a key value pair to redis set

```bash
faas template pull https://github.com/openfaas-incubator/golang-http-template
faas-cli up --build-arg GO111MODULE=on -f add-email.yml'
```