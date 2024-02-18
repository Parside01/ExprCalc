# Expression Calculator

The final project of the second sprint from YandexLMS

## Конфигурация приложения

Прежде чем запускать все в Docker, вы можете поиграть с конфигурацией приложения ./config.yaml

`expression-service` и `app` - часть конфига, которую вы можете взаимодействовать. Она содержит следующие параметры

### app
| Параметр | Тип     | Описание                |
| :------ | :------ | :--------------------- |
| `cache-ttl` | `int` | **Обязательное**. Время в минутах, в течение которого вы хотите хранить выражения в БД |


### expression-service
| Параметр        | Тип     | Описание                |
| :------------- | :------ | :--------------------- |
| `gorutines-count` | `int`    | **Обязательное**. Количество работников, которые будут работать на сервере|
| `expr-queue`     | `string` | **Обязательное**. Имя очереди, в которую будут приходить выражения с frontend|
| `res-queue`      | `string` | **Обязательное**. Имя очереди, в которую будут приходить завершенные задачи|
| `exchange`       | `string` | **Обязательное**. Имя обменника для RabbitMQ|
| `route-key`      | `string` | **Обязательное**. Имя уникального ключа, по которому будут отправляться сообщения в RabbitMQ|
|`worker-info-update`| `int` |  **Обязательное**. Время в секундах, через которое будет обновляться информация о работнике|

## Запуск локально

Клонируйте проект

```bash
git clone https://github.com/Parside01/ExprCalc.git
```

Перейдите в директорию проекта

```bash
cd ExprCalc
```

Запуск приложения в Docker

```bash
docker-compose up --build
```


## После запуска

Нажмите на ссылку http://localhost:80

#### Короткая экскурсия по интерфейсу 

Когда вы заходите в приложение, он будет ожидать вас

![App Screenshot](./screenshots/home-screen.png)

Вы можете ввести математическое выражение в поле ввода

![App Screenshot](./screenshots/example.png)

Перед этим вы можете установить скорость выполнения каждой из поддерживаемых математических операций

![App Screenshot](./screenshots/options.png)


Нажав на кнопку, вы отправляете выражение на сервер, который его обрабатывает, вы можете отслеживать состояние выражения

![App Screenshot](./screenshots/tasks.png)

Вы также можете отслеживать каждого работника в системе

![App Screenshot](./screenshots/worker-monitoring.png)


## Как это работает


### Схема работы

![App Screenshot](./screenshots/scheme.png)


### Стек технологий
#### Фронтенд

<a href="https://reactjs.org/"><img src="https://upload.wikimedia.org/wikipedia/commons/a/a7/React-icon.svg" width="45" style="margin-left:7.5%;"> </a>

#### Бэкенд
<br>
<a href="https://golang.org/"> <img src="https://upload.wikimedia.org/wikipedia/commons/thumb/0/05/Go_Logo_Blue.svg/768px-Go_Logo_Blue.svg.png" width="70" style="margin-left:7%; margin-bottom:15px;"> </a>

<a href="https://github.com/labstack/echo"><img src="https://camo.githubusercontent.com/794ace8f539408352061bb193fce26a0df05bed29d57d2125968fa99143b67cd/68747470733a2f2f63646e2e6c6162737461636b2e636f6d2f696d616765732f6563686f2d6c6f676f2e737667" width="100" style="margin-left:7.5%;margin-bottom:10px;"></a>

<a href="https://www.rabbitmq.com/what-is-rabbitmq.html"><img src="https://storage.yandexcloud.net/media.ref-model.ru/Rabbit_MQ_logo_67175cb0a9.png" width="110" style="margin-left:7.5%; margin-bottom:10px;"><a>

<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/6/64/Logo-redis.svg/330px-Logo-redis.svg.png" width="100" style="margin-left:7.5%;">

 
#### Докер

<a href="https://www.docker.com/"><img src="https://upload.wikimedia.org/wikipedia/commons/thumb/4/4e/Docker_%28container_engine%29_logo.svg/330px-Docker_%28container_engine%29_logo.svg.png" width="120"
style="margin-left:7.5%;"> </a>

## Обратная связь

если у вас есть какие-либо отзывы или проблемы с приложением, пожалуйста, [откройте проблему] (https://github.com/Parside01/ExprCalc/issues/new) или вы можете написать мне в Telegram @parside12

#### Авторы
+ [Parside01](https://github.com/Parside01) :scream_cat:
+ [Kneepy](https://github.com/Kneepy) :japanese_goblin: