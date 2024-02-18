# Expression Calculator

The final project of the second sprint from YandexLMS




## Application Configuration

Before you run everything in Docker, you can play around with the application config ./config.yaml

`expression-service` and `app` - the part of the config that you can interact with. It contains the following fields

### app
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `cache-ttl` | `int` | **Required**. The time in minutes that '\n' expressions will be stored in the DBMS. If it is not necessary, set 0 |


### expression-service
| Parameter        | Type     | Description                |
| :-------------   | :------- | :------------------------- |
| `gorutines-count`| `int`    | **Required**. The number of workers that will be on the server|
| `expr-queue`     | `string` | **Required**. The name of the queue to which expressions from the frontend will arrive |
| `res-queue`      | `string` | **Required**. The name of the queue in which the workers will put the completed tasks| 
| `exchange`       | `string` | **Required**. The name of the exchanger for rabbitmq|
| `route-key`      | `string` | **Required**. The name of the unique key that will be used to send messages to rabbitmq|
|`worker-info-update`| `int` |  **Required**. The time in seconds after which the server will ping the workers|

## Run Locally

Clone the project

```bash
git clone https://github.com/Parside01/ExprCalc.git
```

Go to the project directory

```bash
cd ExprCalc
```

Launching the application in Docker

```bash
docker-compose up --build
```


## After launch

click on the link http://localhost:80

#### Roadmap

When you log in to the app, this will be waiting for you

![App Screenshot](./screenshots/home-screen.png)

You can enter a mathematical expression in the input field

![App Screenshot](./screenshots/example.png)

Before that, you can set the execution speed of each of the supported mathematical operations

![App Screenshot](./screenshots/options.png)

By clicking on the button, you send the expression to the server that processes it, you can monitor the state of the expression

![App Screenshot](./screenshots/tasks.png)

You can also monitor every worker in the system

![App Screenshot](./screenshots/worker-monitoring.png)



## How it works

#### Technology stack


#### Frontend
- HTML
- CSS
- [ReactJS](https://reactjs.org/) <img src="https://upload.wikimedia.org/wikipedia/commons/a/a7/React-icon.svg" width="20">

#### Backend
- [Golang](https://golang.org/) <img src="https://blog.golang.org/go-brand/Go-Logo/SVG/Go-Logo_Aqua.svg" width="20">
- Фреймворк: [Echo](https://github.com/labstack/echo) <img src="https://avatars.githubusercontent.com/u/18666616?s=200&v=4" width="20">

#### Microservices
- Система сообщений: [RabbitMQ](https://www.rabbitmq.com/) <img src="https://www.rabbitmq.com/img/rabbitmq-logo.svg" width="20">

#### Database
- [Redis](https://redis.io/) <img src="https://upload.wikimedia.org/wikipedia/en/6/6b/Redis_Logo.svg" width="20">

#### Docker
- [Docker](https://www.docker.com/) <img src="https://www.docker.com/sites/default/files/d8/2019-07/vertical-logo-monochromatic.png" width="20">





