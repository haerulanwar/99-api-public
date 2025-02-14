# 99-api-public
REST Api for public consume, this is part of 99.co assessment

## How to run

### WARNING
This service cannot run independently as it relies on the user service and listing service. Ensure that both services are running before starting this service.

### Features
This public API provides two main features for retrieving user information for each listing:
1. Fetch user details for each listing by calling the user service API.
   This approach is implemented in the `feature/get-user-by-api` branch.
2. Retrieve user details for each listing from Redis, where the data is initially pooled and subsequently updated whenever there are changes to the user information.
   This approach is implemented in the `main` or `feature/get-user-by-redis-pooling` branch.

These features ensure that each listing includes comprehensive user details, enhancing the overall data quality and usability.


### Required

- redis
- go 1.23.0

### Conf

To set up your environment, use existing config at `.env` or customizing the configurations accordingly.

### Run
```
$ go mod vendor 
$ go run main.go 
```
Now server is running at port 8081 (default port config).

### Public API
#### APIs
##### Get listings
Get all the listings available in the system (sorted in descending order of creation date). Callers can use `page_num` and `page_size` to paginate through all the listings available. Optionally, you can specify a `user_id` to only retrieve listings created by that user.

```
URL: GET /public-api/listings

Parameters:
page_num = int # Default = 1
page_size = int # Default = 10
user_id = str # Optional
```
```json
{
    "result": true,
    "listings": [
        {
            "id": 1,
            "listing_type": "rent",
            "price": 6000,
            "created_at": 1475820997000000,
            "updated_at": 1475820997000000,
            "user": {
                "id": 1,
                "name": "Suresh Subramaniam",
                "created_at": 1475820997000000,
                "updated_at": 1475820997000000,
            },
        }
    ]
}

```

##### Create user
```
URL: POST /public-api/users
Content-Type: application/json
```
```json
Request body: (JSON body)
{
    "name": "Lorel Ipsum"
}
```
```json
Response:
{
    "user": {
        "id": 1,
        "name": "Lorel Ipsum",
        "created_at": 1475820997000000,
        "updated_at": 1475820997000000,
    }
}
```

##### Create listing
```
URL: POST /public-api/listings
Content-Type: application/json
```
```json
Request body: (JSON body)
{
    "user_id": 1,
    "listing_type": "rent",
    "price": 6000
}
```
```json
Response:
{
    "listing": {
        "id": 143,
        "user_id": 1,
        "listing_type": "rent",
        "price": 6000,
        "created_at": 1475820997000000,
        "updated_at": 1475820997000000,
    }
}
```