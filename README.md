## tinyURL

A link shortener service that uses a sequential hashing algorithm to generate a unique short url for a given url and
stores in postgresql database with a in-memory cache.

## Startup

1. Start postgresql server
    ```bash
    docker run --name my-postgres -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres
    bash resources/psql.sh
    ```
2. Make Request:
    ```bash
   curl -X POST -H "Content-Type: application/json" -d '{"url": "https://google.com"}' http://localhost:8080/
   ```

## Todo

1. Config file for DB connection and other parameters
2. Load balancer of some sort
3. Address security concerns of sequential hashing