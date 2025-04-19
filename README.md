# task-management-service

### Problem breakdown and design decisions

The task management service will include basic CRUD functionality. I've chosen to use a NoSQL database due to its
flexibility for adding new fields in the future and its high scalability.

MongoDB, in particular, offers scalability through sharding, allowing us to add new shards to the network as needed.
Additionally, it provides powerful aggregation capabilities for handling complex data queries. Using the ID as the
sharding key will help ensure an even distribution of data across all shards.

### Instructions to run the service

To start the server, update the environment variables in configs/configs.yml and docker.env.dev as per your setup.
Then, run the following command to launch the service on 0.0.0.0:8000:

```makefile
make run_dev
```

To stop the server, use:

```makefile
make kill_dev
```

To run the linter and check code quality:

```makefile
make lint
```

The db documents can be seen using mongo express on the URL http://0.0.0.0:8081/db/taskmanager/task, use the mongo auth
cred present in the docker.env.dev file.

### API Documentation

The service implements simple JWT-based authentication. Clients must include the token in the `Authorization` header
using the format: `Bearer <token>`. The JWT secret key is stored securely in the configuration.
The token payload contains three fields:

- user_id: The user ID from the user service (no whitespace characters)
- user_name: The name of the user
- user_type: The role of the user, either ADMIN or USER

Export the auth token for further use,

```bash
export BEARER_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwidXNlcl9uYW1lIjoiVGVzdCBVc2VyIiwidXNlcl90eXBlIjoiQURNSU4iLCJ1c2VyX2lkIjoidGVzdF91c2VyIiwiaWF0IjoxNTE2MjM5MDIyfQ.FCsWX7t8rP-TAQDSJAhMd-x4xbd-cGDA7mjzySGY7N4
```

#### POST

Creates a new task.

***Request***

```bash
curl --location 'http://0.0.0.0:8000/api/tasks' \
--header "Authorization: Bearer $BEARER_TOKEN" \
--header "Content-Type: application/json" \
--data '{
    "title": "test task - 1",
    "description": "The is a test description."
}'
```

***Response***

```json
{
  "_id": "6804064d52e2de08a3877828",
  "task_id": "6804064d52e2de08a3877828",
  "title": "test task - 1",
  "description": "The is a test description.",
  "owner": "test_user",
  "status": "Created",
  "created_at": "2025-04-19T20:23:41.299290757Z",
  "updated_at": "2025-04-19T20:23:41.299291007Z"
}
```

```bash
export TASK_ID=<task_id from the POST response>
```

#### GET

##### Get a task by ID

Fetches a task using its ID. Only the task owner or an admin is authorized to access it.

***Request***

```bash
curl --location "http://0.0.0.0:8000/api/tasks/$TASK_ID" \
--header "Authorization: Bearer $BEARER_TOKEN"
```

***Response***

```json
{
  "_id": "6804064d52e2de08a3877828",
  "task_id": "6804064d52e2de08a3877828",
  "title": "test task - 1",
  "description": "The is a test description.",
  "owner": "test_user",
  "status": "Created",
  "created_at": "2025-04-19T20:23:41.299Z",
  "updated_at": "2025-04-19T20:23:41.299Z"
}
```

<br>

##### Get all tasks

Returns all tasks if the user is an admin; otherwise, returns only the tasks owned by the user.

- Supports sorting by title, created_at, updated_at, and owner.
    - Default sort: updated_at in descending order.
- Supports filtering by status and owner.
    - Valid status values: Created, InProgress, Succeeded, Discarded.
- Pagination is available via page and recordsPerPage query parameters..

***Request***

```bash
curl --location "http://0.0.0.0:8000/api/tasks?page=1&recordPerPage=4&sortBy=created_at&order=desc&status=Created" \
--header "Authorization: Bearer $BEARER_TOKEN"
```

***Response***

```json
[
  {
    "_id": "6804064d52e2de08a3877828",
    "task_id": "6804064d52e2de08a3877828",
    "title": "test task - 1",
    "description": "The is a test description.",
    "owner": "test_user",
    "status": "Created",
    "created_at": "2025-04-19T20:23:41.299Z",
    "updated_at": "2025-04-19T20:23:41.299Z"
  },
  ...
]
```

#### PUT

Updates an existing task. This replaces the entire task object.
Only admins and task owners are authorized to perform this operation.
Owner of a task cannot be updated.

***Request***

```bash
curl --location --request PUT "http://0.0.0.0:8000/api/tasks/$TASK_ID" \
--header "Authorization: Bearer $BEARER_TOKEN" \
--header 'Content-Type: application/json' \
--data '{
    "title": "test task - 1 (Updated)",
    "status": "InProgress",
    "description": "The description has been updated."
}'
```

***Response***

```json
{
  "_id": "6804064d52e2de08a3877828",
  "task_id": "6804064d52e2de08a3877828",
  "title": "test task - 1 (Updated)",
  "description": "The description has been updated.",
  "owner": "test_user",
  "status": "InProgress",
  "created_at": "2025-04-19T20:23:41.299Z",
  "updated_at": "2025-04-19T20:38:04.644Z"
}
```

#### PATCH

Partially updates a task. Only the fields provided in the request body will be updated.
This action is restricted to the task owner and admin users.

***Request***

```bash
curl --location --request PATCH "http://0.0.0.0:8000/api/tasks/$TASK_ID" \
--header "Authorization: Bearer $BEARER_TOKEN" \
--header 'Content-Type: application/json' \
--data '{
    "description": "The description has been updated twice."
}'
```

***Response***

```json
{
  "_id": "6804064d52e2de08a3877828",
  "task_id": "6804064d52e2de08a3877828",
  "title": "test task - 1 (Updated)",
  "description": "The description has been updated twice.",
  "owner": "test_user",
  "status": "InProgress",
  "created_at": "2025-04-19T20:23:41.299Z",
  "updated_at": "2025-04-19T20:50:50.536Z"
}
```

#### DELETE

Deletes a task. This action is permitted only for the task owner or an admin.

***Request***

```bash
curl --location --request DELETE "http://0.0.0.0:8000/api/tasks/$TASK_ID" \
--header "Authorization: Bearer $BEARER_TOKEN" \
--header 'Content-Type: application/json'
```

***Response***

```json
{
  "status": "Deleted"
}
```

### Scalability

The service can be deployed in any Kubernetes cluster, such as EKS or GKE, with a dedicated MongoDB server.
Since MongoDB is used as the primary database, it offers strong scalability options. By using the `_id` field as the
sharding key, data is evenly distributed across shards, ensuring balanced load and efficient storage.

The server itself is stateless, allowing us to run multiple instances of the service concurrently. Each instance can
independently communicate with the MongoDB backend, making the system horizontally scalable and resilient.

Additionally, Redis can be integrated as a caching layer using a read-through and write-through strategy. This approach
helps to store frequently accessed tasks in memory, reducing database load and improving response times.

### Inter-Service Communication

For inter-service communication between a Task Service and a new User Service, the choice depends on specific needs:

- Use gRPC when low latency, real-time communication is required. It is ideal for request-response patterns where
  immediate feedback is necessary. gRPC uses Protocol Buffers (Protobuf) for efficient serialization, resulting in
  smaller message sizes and faster transmission compared to traditional JSON-based REST APIs.
- Use message queues when asynchronous communication is acceptable. This approach is suitable for scenarios where
  eventual consistency is fine, and you want to decouple services. Message queues enhance fault tolerance, scale well
  with high throughput, and are ideal for background tasks or events.

We can implement both these servers using `Server` interface present at `internal/server` and run the different servers
using command line flags.
