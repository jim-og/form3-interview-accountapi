# form3-interview-accountapi

Jim Gartland

I am new to Go. This assignment has served as a great learning exercise to increase my proficiency with the language and I look forward to programming with it more in the future.

## Usage
Create a new client to access the supported services from the form3 API:

```go
client := form3.NewClient()
```

The client is initialised with the default base url `http://localhost:8080/v1/`. This can be overriden once the client is created:

```go
client := form3.NewClient()
client.BaseURL = "http://localhost:8080/v1/"
```

The context package is used to facilitate the use of deadlines and other services. A valid `context.Context` is a required input parameter when making requests to the client. If you do not wish to configure a custom context, use `context.Background()`:

```go
client := form3.NewClient()
accountID := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"  // a known valid ID
account, resp, err := client.GetAccount(context.Background(), accountID)
```

### Testing
Tests are run with Docker using the provided `docker-compose.yml` from the root directory:

```bash
docker-compose up --build
```

## Design

* Whilst the form3 API appears to always nest request and response data within a `data` key, I've made the decision not to assume that this will always be the case. The impact of this choice is that each client endpoint must wrap the data into a suitable struct with a `json:"data"` field and unwrap it prior to being returned. `accounts::CreateAccount` is an example of this. An additional model called `AccountModel` has been added in `models.go` which specifies the structure an accounts request and response should take. Whilst this requires additional code with each request I feel it is more extendable in the future and keeps common functionality, such as `form3::Do`, agnostic of the form request and response data should take.
* The `go-cmp` package is used in testing to compare the expected output with the actual output. This is sufficient for the cases considered in this assignment but I would expect a production client to have custom comparators for each model to explicitly check fields which require equality and ignore those which are created by the API such as `version`, `createdAt`, etc.
