# Technical Challenge 07 26 2023

A short technical challenge to implement a REST API that satifies the requirements specified in [requirements.txt](./docs/requirements.txt).

## How to Run the API
The project was written using Golang v1.19. Make sure you have installed a compatible version of Golang. Then the API can be started with:
```
go run ./cmd
```

## How to Run Unit Test
Again, the project was written using Golang v1.19. Make sure you have installed a compatible version of Golang. Then the unit tests can be run with:
```
go test -v ./...
```

## Implementation Notes

### Assumptions
 - I am using a path parameter for the DELETE method at /users/{email}.
 - My PATCH method operates more like a PUT method; with users only being an email (their primary key) and a single other field, this feels appropriate.
 - I have made the POST method also require authorization as an admin.
 - I provided thorough unit tests of the API to replace the need for testing using an http client.

### If I Had More Time
 - I would include an external testing suite for the API. Ideally, something "human readable" like Behave tests in Python. I simply did not have the time to research and implement something like this.

### Other
 - The optional open API spec can be found at [openapi.yaml](./api/openapi.yaml)
 - The "hard-coded" user hashes that were used for testing can be found in [hashes.txt](./docs/hashes.txt)