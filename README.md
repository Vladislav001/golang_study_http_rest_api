1. На win10 надо ставить постгрес + руками создать БД (через pgadmin)
(https://winitpro.ru/index.php/2019/10/25/ustanovka-nastrojka-postgresql-v-windows/)


2. Установка cli для миграций:
- скачать https://github.com/golang-migrate/migrate
- закинуть в папку с проектом, напр. сюда C:\Users\Админ\GolandProjects\golang_study_http_rest_api
- исполнить go install .
- проверить, выполнив migrate