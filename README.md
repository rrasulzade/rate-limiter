# Rate Limiter Service

This repository provides an implementation of a Rate Limiter designed to enforce limitations to protect services from excessive requests and protect customers from service interruptions or degradation. By effectively managing access to configured API endpoints, this rate limiter empowers to maintain service reliability.

## How it Works:

The Rate Limiter operates by maintaining a token-based system, where each request consumes a token from a predefined pool. Once the tokens are exhausted within a certain time window, further requests are temporarily restricted, ensuring that resources are allocated efficiently and fairly.

## Get Started:

Explore the repository and follow the guidelines below to leverage rate limiting features for effective traffic management.

### Prerequisites
- Go v1.17
- Docker (optional)

### Installation


#### Using Docker

```bash
# Build the image
docker build -t rate-limiter-srv . 

# Run the container on http://localhost:3003
docker run -p 3003:3003 rate-limiter-srv
```

#### Using local environment

```bash
    # Build the project
    go build -o rate_limiter

    # Run the server
    ./rate_limiter
```

### Testing the server using `curl`
Once the server is up and running, you can simulate client requests.

Use the following command to send a request to the load balancer:

```bash
curl https://localhost:3003/take?endpoint=<RESOURCE_NAME>
```

Replace `<RESOURCE_NAME>` with the name of an API endpoint that is configured in the rate limiter to be restricted. This endpoint represents a specific API endpoint that is subject to rate limiting measures.

## Configuration

The rate limiter is configured using a JSON configuration file located in `config` folder. The configuration includes rate limiting measures for API endpoints.

Sample configuration:
```json
{
  "port": 3003,
  "rateLimitsPerEndpoint": [
    {
      "endpoint": "GET /user/:id",
      "burst": 10,
      "sustained": 6
    },
    {
      "endpoint": "PATCH /user/:id",
      "burst": 10,
      "sustained": 5
    },
    {
      "endpoint": "POST /userinfo",
      "burst": 300,
      "sustained": 100
    }
  ]
}
```
**`port`**: The port number on which the server runs.

**`rateLimitsPerEndpoint`**: An array with a configuration for each API route template the service should provide a rate limit for.

**`endpoint`**: API route template being limited, acts as a key provided by the caller to check and consume its request tokens.

**`burst`**: The number of burst requests allowed.

**`sustained`**: The number of sustained requested per minute.


## Further enhancements
Below are some suggested enhancements to consider for improving the usability of the server in the future.

**Robust logging and monitoring mechanisms**
The server will be equipped with logging and monitoring capabilities to ensure transparency and traceability. 

**Extensive configurations**
Allows for easy customization and management of the server's behavior without the need for code changes. Includes defining various parameters such as database connections, logging levels, API endpoints, and authentication settings.
