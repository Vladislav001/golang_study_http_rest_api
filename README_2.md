1. На win10 надо ставить постгрес + руками создать БД (через pgadmin)
(https://winitpro.ru/index.php/2019/10/25/ustanovka-nastrojka-postgresql-v-windows/)


2. Установка cli для миграций:
- скачать https://github.com/golang-migrate/migrate
- закинуть в папку с проектом, напр. сюда C:\Users\Админ\GolandProjects\golang_study_http_rest_api
- исполнить go install .
- проверить, выполнив migrate
- пример наката миграции: migrate -path C:/Users/Админ/GolandProjects/golang_study_http_rest_api/golang-migrate/cli/migrations -database "postgres://localhost/restapi_dev?sslmode=disable" up
*Но крч дичь какая-то, проблемы, пока через pgadmin ручками тупо через Query Tool вставляю sql


3. Сборка приложения go build -v ./cmd/apiserver
3.2 Запуск  ./apiserver


- Сперва go mod создать go mod init
- Подтянуть пакет go get github.com/stretchr/testify/assert
--------------
TODO:
1. Разобраться с MAKEFILE для win10, мб аналог есть для сборки на win10
2. Разобраться с миграциями, чтобы по норму делались/работали
3. Мб в DOCKER сделать 

4. Сделать логаут
5. Сделать получение списка юзеров (если мы авторизованы)

*Приложить POSTMAN