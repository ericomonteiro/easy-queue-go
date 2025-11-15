# API Documentation (Swagger)

Interactive API documentation powered by OpenAPI/Swagger specification.

## Interactive API Explorer

<!-- swagger-ui: swagger.json -->

## About This Documentation

This interactive documentation is automatically generated from the code comments in the handlers using [swaggo/swag](https://github.com/swaggo/swag).

### Features

- **Try it out**: Test endpoints directly from the browser
- **Authentication**: Support for Bearer token authentication
- **Request/Response examples**: See example payloads for all endpoints
- **Model schemas**: View detailed structure of request and response objects

### How to Use

1. **Explore endpoints**: Click on any endpoint to see details
2. **Try it out**: Click the "Try it out" button to test an endpoint
3. **Authenticate**: For protected endpoints, click the "Authorize" button at the top and enter your Bearer token
4. **Execute**: Fill in the required parameters and click "Execute" to make a real API call

### Authentication

Protected endpoints require a Bearer token. To authenticate:

1. Login using the `/auth/login` endpoint
2. Copy the `access_token` from the response
3. Click the **Authorize** button (🔓) at the top of the page
4. Enter: `Bearer {your_access_token}`
5. Click **Authorize**

Now you can test protected endpoints!

## Generating Documentation

The Swagger documentation is generated from code comments. To update:

```bash
make swagger-generate
```

See [Swagger Documentation Guide](api/swagger.md) for more details on how to add/update API documentation.
