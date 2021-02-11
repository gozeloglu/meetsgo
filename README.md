# MeetsGo

## What is MeetsGo?

MeetsGo comes from "Let's Meetup". I have written in Go, so that I added "Go" phrase as postfix.

MeetsGo is a REST API written in Go. It can be seen as basic REST API which handles the CRUD operations for Meetup application. I started to learn Go, so that I mainly focused on backend and CRUD operations, not frontend or advance topics. 

## Tech Stack

* **Go-1.15.7**
* **Gorm**
* **net/http**
* **PostgreSQL**
* **gorilla/mux**

## How to run?

You need to install Go. You can follow the instructions [here](https://golang.org/dl/) to download and install the Go. Then, you need to install Go packages.

You can install [`Gorm`](https://gorm.io/) packages with the following commands. 

```shell
$ go get -u gorm.io/gorm
$ go get -u gorm.io/driver/postgres   # install postgresql driver 
```

You can install [`gorilla/mux`](https://github.com/gorilla/mux) package with the following command.

```shell
$ go get -u github.com/gorilla/mux
```

To run the project, type the following command:

```shell
$ go run .
```

To run the test, type the following command:

```shell
$ go test
```

## Database Schema

Database table structure is simple. There are two different table with different attributes.


### User Table
| id | username | name | surname | password | email | age | is_admin |
| ---| -------- | ---- | ------- | -------- | ----- | --- | -------- |
| 1  | john  | John   | Wick |$2a$10$PomZY2/t5lXSHjAx4J5bmexHnx3C6fo3DYF5Yp08CFNMKQmD9l6LW  | john@mail.com | 34 | false |
| 2  | ahmet | Ahmet | Yilmaz | $2a$10$0FjVe/kKu833DpF7skkVzOQnqWMxTKyrt.3In/9xgY.Bb.xA6Fcbe | ahmet@email.com | 32 |true | 
| ... | ... | ... | ... | ... | ... | ... | ... | 
### Meetup Table
| id | meetup_name | meetup_details | start_date | end_date | address | quota | registered_user_count  |
| ---| -------- | ---- | ------- | -------- | ----- | --- | -------- |
| 1  | GDG Ankara GoLang  | We are going to learn Go on live session. | 02.04.2021T12:00:00 | 02.04.2021T15:00:00| YouTube Channel | 100 | 13 |
| 2  | GDG Ankara Flutter  | We are going to learn Flutter on live session. | 03.04.2021T12:00:00 | 03.04.2021T15:00:00| YouTube Channel | 100 | 7 | 
| ... | ... | ... | ... | ... | ... | ... | ... |
 
## API

### User

* Create User - `/user/create` - POST
* Get User - `/user/{username}` - GET
* Get Users - `/users` - GET 
* Login - `/user/login` - POST
* Update User Profile - `/user/update/{username}` - PUT

### Meetup

* Create Meetup - `/meetup/create/{admin_username}` - POST
* Get Meetups - `/meetups` - GET
* Get Meetup Details - `/meetup/details/{meetup_id}` - GET
* Delete Meetup - `/meetup/delete/{meetup_id}` - DELETE
