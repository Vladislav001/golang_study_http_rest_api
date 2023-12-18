Проект представляет собой тестовый сервер, способный обработать регистрацию, авторизацию, получение профиля

1. Сборка приложения go build -v ./cmd/apiserver
2. Запуск  ./apiserver
3. POSTMAN в https://github.com/Vladislav001/golang_study_http_rest_api/blob/5666219ce6be9e656990dfdbe767057b1b7246b2/golang_study_http_rest_api.postman_collection.json
4. Запуск тестов: cd internal/app/store/teststore и затем go test - для запуска в конкретной папки. Если надо тесты всего проекта, то  go test ./... из корня проекта