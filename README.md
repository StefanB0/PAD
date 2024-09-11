# Web Proxy: Image sharing for novice artists

- Author: Boicu Stefan
- Academic group: FAF203

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

## Deployment and Scaling

- **Docker**. All the services will be dockerised for deployment to ensure they behave in a predictable manner.
- **Docker compose**. Scaling will be achieved by using docker compose because of its simplicity.

**Communication**

- All servers use HTTP for communication
- The gateway and service discovery use REST API
- The image service, user service and recommendation service use RPC over HTTP

## Technology Stack
- **User service:** Go
  - **User database:** Postresql
- **Image service:** Go
  - **Image database:** MongoDB
- **Recommendation service:** Go
  - **Analytics database:** Postresql
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
- Status

#### Status

All services have a status endpoint that returns status code 200 and message OK.

- **Endpoint:** GET `/status`
- **Request:** None
- **Response:** `OK`

### Gateway

### Get Image

Retrieves the image content based on the provided ID.

- **Endpoint:** GET `/image/:id`
- **Request:** None
- **Response:**
  - `200 OK` with image content (JPEG)
  - `404 Not Found` if the image is not available

### Get Image Info

Retrieves information about the image based on the provided ID.

- **Endpoint:** GET `/image/info/:id`
- **Request:** None
- **Response:**
  - `200 OK` with JSON containing image information
  - `404 Not Found` if the image is not available

### Upload Image

Uploads a new image with the provided parameters.

- **Endpoint:** POST `/image`
- **Request:** Params with "token," "author," "title," "description," "tags," and "image" fields
- **Response:**
  - `201 Created` with JSON containing the uploaded image details
  - Status code indicating the error if unsuccessful

### Like Image

Likes a specified image.

- **Endpoint:** POST `/image/:id/like`
- **Request:** Params with "id" field
- **Response:**
  - `200 OK` with JSON containing information about the liked image
  - Status code indicating the error if unsuccessful

### Delete Image

Deletes a specified image.

- **Endpoint:** DELETE `/image/:id`
- **Request:** Params with "id" field
- **Response:**
  - `200 OK` with JSON containing information about the deleted image
  - Status code indicating the error if unsuccessful

### Update Image

Updates information for a specified image.

- **Endpoint:** PUT `/image/:id`
- **Request:** Params with "id," "author," "title," and "description" fields
- **Response:**
  - `200 OK` with JSON containing information about the updated image
  - Status code indicating the error if unsuccessful

### Get Recommendations

Retrieves recommendations based on a specified tag.

- **Endpoint:** GET `/recommend/:tag`
- **Request:** Params with "tag" field
- **Response:**
  - `200 OK` with JSON containing recommended image ID
  - `404 Not Found` if no recommendations for the given tag

### Get Tags

Retrieves a list of tags.

- **Endpoint:** GET `/tags`
- **Request:** None
- **Response:**
  - `200 OK` with JSON containing an array of tags
  - `404 Not Found` if the tags are not available

### User Service

#### Register User

#### Login User

### Image Service

### Get Image

Retrieves the image content based on the provided ImageID.

- **Endpoint:** POST `/getImage`
- **Request:** `{ "imageID": 123 }`
- **Response:** `200 OK`
  - Body: Image content (JPEG)

### Get Image Info

Retrieves information about the image based on the provided ImageID.

- **Endpoint:** POST `/getImageInfo`
- **Request:** `{ "imageID": 123 }`
- **Response:** `200 OK`
  - Body: JSON containing image information (ImageID, Author, Title, Description, Tags)

### Upload Image

Allows users to upload images to the server.

- **Endpoint:** POST `/uploadImage`
- **Request:** Multipart/form-data with image file and metadata
  - author: string
  - title: string
  - description: string
  - tags: string
  - image: file with .jpg extension
- **Response:** `201 Created`
  - Body: JSON containing the newly created image ID

### Like Image

Increments the like count for the specified image.

- **Endpoint:** POST `/likeImage`
- **Request:** `{ "imageID": 123 }`
- **Response:** `200 OK`
  - Body: "Image liked"

### Update Image

Updates information for the specified image.

- **Endpoint:** POST `/updateImage`
- **Request:** accepts multiform or json
  ```
  {
    token: "string",
    imageID: 0
    author: "string"
    title: "string"
    description: "string"
  }
  ```
- **Response:** `200 OK`
  - Body: "Image updated"

### Delete Image

Deletes the specified image.

- **Endpoint:** POST `/deleteImage`
- **Request:** JSON with token and imageID
- **Response:** `200 OK`
  - Body: "Image deleted"


### Service Discovery

#### Get All Services

Retrieves information about all services.

- **Endpoint:** GET `/serviceall`
- **Request:** None
- **Response:** `200 OK`
  - Body: JSON containing information about all services

#### Get Service by Name

Retrieves information about a specific service by name.

- **Endpoint:** GET `/service/:name`
- **Request:** None
- **Response:** 
  - `200 OK` with JSON containing information about the service
  - `404 Not Found` if the service is invalid

#### Add Service

Adds a new service with the provided name and address.

- **Endpoint:** POST `/service`
- **Request:** JSON with "name" and "address" fields
- **Response:**
  - `201 Created` with JSON containing the secret key
  - `400 Bad Request` if the request is missing required fields

#### Remove Service

Removes a service by name.

- **Endpoint:** DELETE `/service/:name`
- **Request:** Params with "name," "address," and "secretkey"
- **Response:**
  - `200 OK` if the service is successfully removed
  - `401 Unauthorized` if removal is unauthorized

### Recommendation service

#### Get Tags

Retrieves a list of tags.

- **Endpoint:** POST `/getTags`
- **Request:** None
- **Response:** `200 OK`
  - Body: JSON containing an array of tags

#### Get Recommendations

Retrieves recommendations based on a specified tag.

- **Endpoint:** POST `/getRecommendations`
- **Request:** JSON with "tag" field
- **Response:** 
  - `200 OK` with JSON containing recommended image ID
  - `404 Not Found` if there are no recommendations for the given tag

#### Add Image

Adds a new image with the provided ID and tags.

- **Endpoint:** POST `/addImage`
- **Request:** JSON with "id" and "tags" fields
- **Response:** `201 Created`
  - Body: "Image added"

#### Update Image

Updates the views and likes for a specified image.

- **Endpoint:** POST `/updateImage`
- **Request:** JSON with "id," "views," and "likes" fields
- **Response:** `200 OK`
  - Body: "Image updated"

#### Delete All Data

Deletes all data.

- **Endpoint:** POST `/deleteALL`
- **Request:** None
- **Response:** `200 OK`
  - Body: "All data deleted"