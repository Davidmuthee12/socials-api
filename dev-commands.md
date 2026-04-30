# Dev Commands

## Docker

- List containers:
  `docker ps`

## Redis cache checks

- Show all keys:
  `docker exec -it learnhouse-redis-dev redis-cli KEYS "*"`
- Show user keys:
  `docker exec -it learnhouse-redis-dev redis-cli KEYS "user-*"`
- Read cached user value:
  `docker exec -it learnhouse-redis-dev redis-cli GET user-119`

Note: `KEYS user-119` checks key names by pattern. Use `GET user-119` to read the stored JSON value.

## Load testing (autocannon)

- Run load test against user endpoint:
  `npx autocannon "http://localhost:8080/v1/users/119/" --connections 10 --duration 5 -H "Authorization: Bearer <TOKEN>"`

- Run health endpoint load test at 25 req/s for 1 second:
  `npx autocannon -r 25 -d 1 -c 1 --renderStatusCodes http://localhost:8080/v1/health`

- Run health endpoint load test at 20 req/s for 1 second:
  `npx autocannon -r 20 -d 1 -c 1 --renderStatusCodes http://localhost:8080/v1/health`

- Run health endpoint load test at 40 req/s for 2 seconds with 10 connections:
  `npx autocannon -r 40 -d 2 -c 10 --renderStatusCodes http://localhost:8080/v1/health`

Note: use `--connections` and `--duration` (not `==connections` or `==duration`).

Note: `/v1/health` is public now, so successful requests return `200`.

Note: the global rate limiter is still enabled with a limit of 20 requests per 5-second window, so health endpoint load tests will start returning `429 Too Many Requests` once that window is exceeded.

Observed results:

- `npx autocannon -r 25 -d 1 -c 1 --renderStatusCodes http://localhost:8080/v1/health`
  produced `200 x 20` and `429 x 13`.
- `npx autocannon -r 20 -d 1 -c 1 --renderStatusCodes http://localhost:8080/v1/health`
  produced `200 x 20` and `429 x 20`.
- `npx autocannon -r 40 -d 2 -c 10 --renderStatusCodes http://localhost:8080/v1/health`
  produced `200 x 200` and `429 x 663`.

## Check number of lines in entire codebase

`Get-ChildItem -Recurse -Include *.go | Get-Content | Measure-Object -Line`
