# Expression Calculator

The final project of the second sprint from YandexLMS


## Run Locally

Clone the project

```bash
    git clone https://github.com/Parside01/ExprCalc.git
```

Go to the project directory

```bash
    cd ExprCalcS
```

Launching the application in Docker

```bash
    docker-compose up --build
```


It is worth noting that it is better to launch in Docker, if you do not want to do this through Docker, then you will have to change the standard application config. 


## Application Configuration

Before you run everything in Docker, you can play around with the application config ./config.yaml

`expression-service` and `app` - the part of the config that you can interact with. It contains the following fields

### app
| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `cache-ttl` | `int` | **Required**. The time in minutes that '\n' expressions will be stored in the DBMS. If it is not necessary, set 0 |
