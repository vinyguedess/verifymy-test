![UnitTest](https://github.com/vinyguedess/verifymy-test/actions/workflows/unit_tests.yml/badge.svg)

# VerifyMy: Test
Test API for VerifyMy interview process.

## Getting started
First of all you need to check environment variables that are loaded from a `.env` file which is not versioned for security purposes. There's a file called `.env.default` which you can copy into `.env` and change its values if you want to run the application locally.

To ease everything, project is executed in a golang container. To run it, you need to have docker installed in your machine.

### Lifting up the container
To lift up the container, run the following command:
```bash
make up
```

It will lift up `app` and `mysql` container and ssh you into it.

### Running tests
To run tests, run the following command inside the container:
```bash
make test
```

It will generate required mocks, run tests and generate coverage report inside `coverage` folder.

### Running application
To run the application, run the following command inside the container:
```bash
make dev
```

It will run the application in `6073` port.

## Documentation
API documentation was done with Swagger. To access the it, run the application and access `http://localhost:6073/swagger/`. All endpoints are documented there, except for healtcheck endpoint that is a `GET /` where you can check
service name and its version.
