![UnitTest](https://github.com/vinyguedess/verifymy-test/actions/workflows/unit_tests.yml/badge.svg)

# VerifyMy: Test
Test API for VerifyMy interview process.

## Getting started
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

## Endpoints

### Sign in
**POST** `/auth/sign_in`

#### Request
| Attribute | Type | Description |
| --- | --- | --- |
| name | string | User name |
| date_of_birth | string | User date of birth (YYYY-MM-DD) |
| email | string | User email |
| password | string | User password |
| address | string | User address |

#### Response
| Attribute | Type | Description |
| --- | --- | --- |
| access_token | string | Token used for private routes |
| expires_at | integer | Token expiration date |
