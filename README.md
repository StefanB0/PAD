# Laboratory work nr 1

# Web Proxy: Image sharing for novice artists

- Author: Boicu Stefan
- Academic group: FAF203

# Lab 2 goals
0. ***Draw new diagrams***
1. Just be
2. Trip circuit breaker
3. Service high availability
   - Solution: redundant servers
4. Graphana + Prometheus
5. 2 Phase commits. Affected endpoint: post image
6. Consistent Hashing for Cache
   - Cache image requests in several redis chaches
   - Implement Consistent caching for these redis instances
7. Cache High availability 
   - Redundancies for the redis chaches
8. Saga pattern
9. Database redundancy
   - Replicatoin of the image mongodb database.

# Design Document

## Running the lab

- Open the terminal in the root folder of the project
- To build all the docker containers run `make build-all`
- To pull all images from dockerhub run `make pull-all`
- To run all the containers run `docker compose up -d`
- The main access point to the system is the gateway which is revealed at port 8080

## Architecture choices
- Service high availability will be achieved by running multiple servers and redirecting if a request failed
- To aggregate data, Prometheus + Grafana will be used.
- The saga pattern will be used for the image creation endpoint.
- For consistent hashing, several new redis instances will be created for image requests.
- Cache high availability will be achieved by sharing hashes between redis instances, with the system still working if half the instances go down
- Database redundancy will be applied to the mongodb image database by using replication.

## Service Boundaries
- **User service:** The user service assures user authorisation. It lets users register and delivers timed access tokens on login.
- **Image service:** The image service handles operations related to images. Depending on user authorisation, images can be uploaded and deleted. Images can have an optional name and short description.
- **Recommendation Service:** The feed service will compile the user's interests and give him an image (or set of images) based on their preferences.
- **API Gateway:** The API Gateway will stand between the (supposed) front end of the website and the other services. It will query service discovery and then forward the requests to the user and image service.
- **Service Discovery:** The service discovery will keep the addresses of all other services in memory and be queried by services which make requests. It will also perform load balancing.
- **Cache:** The cache will store image requests for a couple of minutes in case the same image is requested again.

![Design Diagram](local/image/PAD_LAB_1.jpg)

**Communication**

- All servers use HTTP for communication
- The gateway and service discovery use REST API
- The image service, user service and recommendation service use RPC over HTTP

## Technology Stack
- **User service:** Go
  - **User database:** MongoDB
- **Image service:** Go
  - **Image database:** MongoDB
- **API Gateway:** Ruby on Rails.
- **Service Discovery:** Ruby
- **Cache:** Redis db for cache.
- **Prometheus:** Monitor the databases and servers, compile statistics
- **Graphana:** Pull data from prometheus and display it 


## Data Management

### Operations

- Register user
- Login user
- Get image
- Get image info
- Upload image
- Delete image
- Modify image name or description
- Get tags
- Get recommended image

### User Service

#### Register User

- gRPC request

``

```json
{
  "name":     "Mario",
  "password": "Mario"
}

```

- gRPC response

```json
{
  "message":  "success",
  "error":    "error message"
}

```

On success the message is "success" and the error is an empty string. On failure the message is "failure" and error messages vary.

Errors will be "empty/invalid name and/or password" as well as "existing username"

#### Login User

- gRPC request

```json
{
  "name":     "Mario",
  "password": "Mario"
}

```

- gRPC response

```json
{
  "accessToken":  "abcdnndsnodnsoknfodnsfdsnflsdknfsldkn",
  "refreshToken": "fghsadsladnsknvuoamdoxlacmnndsaiodajj",
  "error":        "error message"
}

```

On successful request the user is returned an access and refresh token. The error message is empty.

On failure the user gets an error message. "Your login or password is not correct"

### Image Service

#### Upload image

- gRPC request

```json
{
  "imageBlob":    "0001010101010101010101010101",
  "author":       "gilgamesh777",
  "title":        "post-post-modern Mona Lisa", // optional
  "description":  "A new reimagining of a classic painting" //optional
}

```

The requests consist of an image blob and optionally a title and a description

- gRPC response

```json
{
  "message":  "success",
  "imageID":  "ffnn990",
  "error":    "error message"
}

```

On success the message is success and the error is empty. On failure the message is failure and the image ID is empty. Error contains error message.

#### Get image

- gRPC request

```json
{
  "imageID": "ffnn990"
}

```

- gRPC response

```json
{
  "imageBlob":    "0001010101010101010101010101",
  "author":       "gilgamesh777",
  "title":        "post-post-modern Mona Lisa", // if present
  "description":  "A new reimagining of a classic painting", // if present
  "error":        "error message"
}

```

In case the image is present in the database, the image blob is returned. If there is no such image the image blob is empty and an error message is returned.

#### Modify image

- gRPC request

```json
{
  "imageID":      "ffnn990",
  "title":        "post-post-modern Mona Lisa", // optional
  "description":  "A new reimagining of a classic painting" // optional
}

```

At least one between the title and the description should be present or it returns an error

- gRPC response

```json
{
  "message":  "success",
  "error":    "error message"
}

```

If the operation is successful, the "success" message is returned, otherwise an error is returned.

#### Delete image

- gRPC request

```json
{
  "imageID":      "ffnn990"
}

```

- gRPC response

```json
{
  "message":  "success",
  "error":    "error message"
}

```

If the operation is successful, the "success" message is returned, otherwise an error is returned.

### API Gateway

The gateway mirrors the endpoints from the image and user service, but it also checks the redis cache to see if the user is authorised to perform the actions it wants to.

#### Register User

- http request
- endpoint: `/signup`
- method: POST

```json
{
  "name":     "Mario",
  "password": "Mario"
}

```

- http response

```json
{
  "message":  "success",
  "error":    "error message"
}

```

On success the message is "success" and the error is an empty string. On failure the message is "failure" and error messages vary.

Errors will be "empty/invalid name and/or password" as well as "existing username"

#### Login User

- http request
- endpoint: `/login`
- method: POST

```json
{
  "name":     "Mario",
  "password": "Mario"
}

```

- http response

```json
{
  "accessToken":  "abcdnndsnodnsoknfodnsfdsnflsdknfsldkn", // as cookie
  "refreshToken": "fghsadsladnsknvuoamdoxlacmnndsaiodajj", // as cookie
  "error":        "error message"
}

```

On successful request the user is returned an access and refresh token. The error message is empty.

On failure the user gets an error message. "Your login or password is not correct"

#### Refresh token

#### Upload image

- http request
- endpoint: `/img`
- method: POST

```json
{
  "accessToken":  "abcdnndsnodnsoknfodnsfdsnflsdknfsldkn",
  "imageBlob":    "0001010101010101010101010101",
  "title":        "post-post-modern Mona Lisa", // optional
  "description":  "A new reimagining of a classic painting" //optional
}

```

- http response

```json
{
  "message":  "success",
  "imageID":  "ffnn990",
  "error":    "error message"
}

```

#### Get image

- http request
- endpoint: `/{imgID}`
- method: GET

```json
{}

```

- http response

```json
{
  "imageBlob":    "0001010101010101010101010101",
  "author":       "gilgamesh777",
  "title":        "post-post-modern Mona Lisa", // if present
  "description":  "A new reimagining of a classic painting", // if present
  "error":        "error message"
}

```

In case the image is present in the database, the image blob is returned. If there is no such image the image blob is empty and an error message is returned.

#### Modify image

- http request
- endpoint: `/{imgID}`
- method: PATCH

```json
{
  "accessToken":  "abcdnndsnodnsoknfodnsfdsnflsdknfsldkn",
  "imageID":      "ffnn990",
  "title":        "post-post-modern Mona Lisa", // optional
  "description":  "A new reimagining of a classic painting" // optional
}

```

At least one between the title and the description should be present or it returns an error

- http response

```json
{
  "message":  "success",
  "error":    "error message"
}

```

If the operation is successful, the "success" message is returned, otherwise an error is returned.

#### Delete image

- http request
- endpoint: `/{imgID}`
- method: DELETE

```json
{
  "accessToken":  "abcdnndsnodnsoknfodnsfdsnflsdknfsldkn",
  "imageID":      "ffnn990"
}

```

- http response

```json
{
  "message":  "success",
  "error":    "error message"
}

```

### Service Discovery

#### Get Service

- http request
- endpoint: `/service`
- method: GET

```json
{
  "serviceName": "USER SERVICE"
}

```

- http response

```json
{
  "host":       "localhost:3799",
  "connTicket": "B17A90JKL"
}

```

The user service is only available on the private network. A service asks for what it wants to acces and the service discovery returns the address and a ticket, while also performing load balancing.

After closing a connection, the respective service should close the connection ticket

#### Close ticket

- http request
- endpoint: `/ticket`
- method: POST

```json
{
  "connTicket": "B17A90JKL"
}

```

- http response

```json
{
  "message":  "success",
  "error":    "error message"
}

```

#### Add service

- http request
- endpoint: `/service`
- method: POST

```json
{
  "serviceName": "USER SERVICE",
  "host": "localhost:3799"
}

```

- http response

```json
{
  "message":  "success",
  "error":    "error message"
}

```

#### Remove service

- http request
- endpoint: `/service`
- method: DELETE

```json
{
  "serviceName": "USER SERVICE",
  "host": "localhost:3799"
}

```

- http response

```json
{
  "message":  "success",
  "error":    "error message"
}

```

This method is used when a service performs a graceful shutdown. If a service crashes the Service Discovery will detect it using the heartbeat algorithm.

## Deployment and Scaling

- **Docker**. All the services will be dockerised for deployment to ensure they behave in a predictable manner.
- **Docker compose**. Scaling will be achieved by using docker compose because of its simplicity.
