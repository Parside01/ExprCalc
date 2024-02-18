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
- [ReactJS](https://reactjs.org/) <br>

<a href="https://reactjs.org/"><img src="https://upload.wikimedia.org/wikipedia/commons/a/a7/React-icon.svg" width="45" style="margin-left:20px;"> </a>

#### Backend

- [Golang](https://golang.org/): <br>
  
<a href="https://golang.org/"> <img src="https://upload.wikimedia.org/wikipedia/commons/thumb/0/05/Go_Logo_Blue.svg/768px-Go_Logo_Blue.svg.png" width="70" style="margin-left:10px;"> </a>

- [Go-echo](https://github.com/labstack/echo): <br> 

<a href="https://github.com/labstack/echo"><img src="https://camo.githubusercontent.com/794ace8f539408352061bb193fce26a0df05bed29d57d2125968fa99143b67cd/68747470733a2f2f63646e2e6c6162737461636b2e636f6d2f696d616765732f6563686f2d6c6f676f2e737667" width="100" style="margin-left:20px;"></a>

- [RabbitMQ](https://www.rabbitmq.com/) <br>

<a href="https://www.rabbitmq.com/what-is-rabbitmq.html"><img src="https://storage.yandexcloud.net/media.ref-model.ru/Rabbit_MQ_logo_67175cb0a9.png" width="110" style="margin-left:25px;"><a>

- [Redis](https://redis.io/) <br> 

<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/6/64/Logo-redis.svg/330px-Logo-redis.svg.png" width="100" style="margin-left:25px;">
 
#### Docker
- [Docker](https://www.docker.com/) <br>

<a href="https://www.docker.com/"><img src="https://upload.wikimedia.org/wikipedia/commons/thumb/4/4e/Docker_%28container_engine%29_logo.svg/330px-Docker_%28container_engine%29_logo.svg.png" width="120"
style="margin-left:25px;"> </a>





