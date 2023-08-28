<a name="readme-top"></a>


![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/123DaNIS123/UsersSegmentsAPI)


<!-- PROJECT LOGO -->
<div align="center">

  <h3 align="center">UsersSegments Golang REST API</h3>

  <p align="center">
    An awesome REST API on Golang with Gin, Gorm, PostgreSQL and Swagger!
    <br />
    <a href="#readme-docs"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/123DaNIS123/UsersSegmentsAPI/issues">Request Feature</a>
    ·
    <a href="https://github.com/123DaNIS123/UsersSegmentsAPI/issues">Report Bug</a>
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">With Docker</a></li>
        <li><a href="#installation">Without Docker</a></li>
      </ul>
    </li>
    <li><a href="#usage">Documentation</a></li>
    <li><a href="#roadmap">Features</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project


Pretty good REST API for Users Segments service. You can add users and segments. You can connect users to segments. You can add users to created segments automatically by writing the percentage of users you need and etc.




### Built With

REST API built on Golang with Gin framework and Gorm ORM, docs added with Swagger. PostgreSQL DB.


* [![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/123DaNIS123/UsersSegmentsAPI)](https://go.dev/)
* [Gin Framework](https://gin-gonic.com/)
* [GORM ORM](https://gorm.io/)
* [Swagger](https://github.com/swaggo/gin-swagger)
* [PostgreSQL](https://www.postgresql.org/)




<!-- GETTING STARTED -->
## Getting Started

How to run the project.

Watch .env file for environment variables. Change them if you need

### With Docker

Write this in the folder where `docker-compose.yml` is located

  ```sh
  docker-compose up --build
  ```

### Without Docker
Write this in the folder where `main.go` is located
1. Go build
   ```sh
   go build
   ```
2. Go run `main.go`
   ```sh
   go run `main.go`
   ```

<a name="readme-docs"></a>
## Documentation:
You can watch whole documentation in swagger in: `hostname:port/docs/index.html` when running the project

### Necessary API Docs: 
    POST /segment

    Description: Creates segmetn with given name and percentage

    parameters: json:
    {
        "name": string,
        "percentage": int // not mandatory
    }
____
    DELETE /segment

    Description: Deletes segment with given name

    parameters: json:
    {
        "name": string,
    }
____
    POST /bind

    Description: Add user with user_id to segments_add segments and removes from segments_remove segments. 

    parameters: json:
    {
        "segments_add": [string], //not mandatory
        "segments_remove": [string], //not mandatory
        "user_id": int
    }
____
    POST /userbinds

    Description: Shows all segments user with user_id attached to

    parameters: json:
    {
        "user_id": int
    }
____
    POST /timedata

    Description: Get all segment binding actions since given year and month

    parameters: json:
    {
        "year": int,
        "month": int
    }


<!-- ROADMAP -->
## Features:

- [x] Метод создания сегмента. Принимает slug (название) сегмента.
- [x] Метод удаления сегмента. Принимает slug (название) сегмента.
- [x] Метод добавления пользователя в сегмент. Принимает список slug (названий) сегментов которые нужно добавить пользователю, список slug (названий) сегментов которые нужно удалить у пользователя, id пользователя.
- [x] Метод получения активных сегментов пользователя. Принимает на вход id пользователя.
- [x] ☐ Покрытие кода тестами - частично
- [x] Swagger файл для вашего API.
  ____
- [x] Доп. задание 1:

    Иногда пользователи приходят в поддержку и спрашивают почему у них пропал/появился какой-то новый функционал. Нужно иметь возможность посмотреть когда точно пользователь попал в конкретный сегмент.

    Задача: реализовать сохранение истории попадания/выбывания пользователя из сегмента с возможностью получения отчета по пользователю за определенный период. На вход: год-месяц. На выходе ссылка на CSV файл.

    Пример отчета:

    идентификатор пользователя 1;сегмент1;операция (добавление/удаление);дата и время

    идентификатор пользователя 1;сегмент2;операция (добавление/удаление);дата и время

    идентификатор пользователя 2;сегмент3;операция (добавление/удаление);дата и время
    ____

- ☐ Доп. задание 2 - В процессе:

    Бывают ситуации когда нам нужно добавить пользователя в эксперимент на ограниченный срок. Например выдать скидку всего на 2 дня.

    Задача: реализовать возможность задавать TTL (время автоматического удаления пользователя из сегмента)

    Пример: Хотим чтобы пользователь попал в сегмент на 2 дня - для этого в метод добавления сегментов пользователю передаём время удаления пользователя из сегмента отдельным полем
    _____

- [x] Доп. задание 3:

    Мы хотим добавлять пользователя в сегмент не в ручную, а автоматически. В сегмент будет попадать заданный процент пользователей.

    Задача: в методе создания сегмента, добавить опцию указания процента пользователей, которые будут попадать в сегмент автоматически. В методе получения сегментов пользователя, добавленный сегмент должен отдаваться у заданного процента пользователей.

    Пример: создали сегмент AVITO_VOICE_MESSAGES и указали что 10% пользователей будут попадать в него автоматически. Пользователь 1000 попал в этот сегмент автоматически. При запросе сегментов пользователя 1000, сегмент AVITO_VOICE_MESSAGES должен отдаваться всегда.
    ____





<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>


## Contact

Danis Nizamutdinov
<br>My Telegram: @DDDaNiSSS
