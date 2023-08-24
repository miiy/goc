## swagger-ui

<https://github.com/swagger-api/swagger-ui>

<https://github.com/swagger-api/swagger-ui/blob/master/docs/usage/configuration.md>

```bash
docker pull swaggerapi/swagger-ui
docker run --name swagger-ui -p 8080:8080 -e URLS=[ { url: 'https://petstore.swagger.io/v2/swagger.json', name: 'Petstore' } ] -v /data/www/swagger.zsaix.com:/app -d swaggerapi/swagger-ui
```
