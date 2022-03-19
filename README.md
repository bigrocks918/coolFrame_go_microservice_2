# Working with Microservices in Go

This is the source code for the Udemy course **Working with Microservices and Go**. This project
consists of a number of loosely coupled microservices, all written in Go:

- broker-service: an optional single entry point to connect to all services from one place (accepts JSON)
- authentication-service: authenticates users against a Postgres database (accepts JSON)
- logger-service: logs important events to a MongoDB database (accepts RPC and JSON)
- queue-listener-service: consumes messages from amqp (RabbitMQ) and initiates actions based on payload
- mail-service: sends email (accepts JSON)

All services (except the broker) register their access urls with etcd, and renew their leases automatically.
This allows us to implement a simple service discovery system, where all service URLs are accessible with
"service maps" in the Config type used to share application configuration in the broker service.

In addition to those services, the included `docker-compose.yml` at the root level of the project
starts the following services:

- Postgresql 14
- mailhog
- MongoDB

## Running the project
From the root level of the project, execute this command (this assumes that you have 
[GNU make](https://www.gnu.org/software/make/) and a recent version
of [Docker](https://www.docker.com/products/docker-desktop) installed on your machine):

~~~
make up
~~~

Then start the front end:

~~~
make start
~~~


Hit the front end with your web browser at `http://localhost:80`

To stop everything:

~~~
make stop
make down
~~~

All make commands:

~~~
tcs@Grendel go-microservices % make help
 Choose a command:
  up         starts all containers in the background without forcing build
  up_build   stops docker-compose (if running), builds all projects and starts docker compose
  auth       stops authentication-service, removes docker image, builds service, and starts it
  broker     stops broker-service, removes docker image, builds service, and starts it
  logger     stops logger-service, removes docker image, builds service, and starts it
  mail       stops mail-service, removes docker image, builds service, and starts it
  listener   stops listener-service, removes docker image, builds service, and starts it
  down       stop docker compose
  start      starts the front end
  stop       stop the front end
  test       runs all tests
  clean      runs go clean and deletes binaries
  help       displays help
~~~