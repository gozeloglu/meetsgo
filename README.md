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


## API

* Create User - `/user/create` - POST
* Get User - `/user/{username}` - GET
* Get Users - `/users` - GET 
* Login - `/user/login` - POST
* Update User Profile - `/user/update/{username}` - PUT
