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

Note: use `--connections` and `--duration` (not `==connections` or `==duration`).
