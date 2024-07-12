# weather-app
Приложение "Погода" предоставляет API для работы с информацией о погоде.

## Что сделано

Часть 1 (получение данных): 
1. Получение данных о городах с помощью geocoding-api.
2. Получение данных о погоде в города с помощью open weather map.
3. Данные о погоде обновляются асинхронно в фоновом процессе.
4. Реализована возможность параллельного извлечения данных.
5. Описаны файлы миграции.

Часть 2 (API-ручки): 
1. Получение списка городов, для которых доступен прогноз.
2. Краткий прогноз для города.
3. Детальный прогноз для города на конкретную дату/время.
4. Регистрация пользователей.
5. Авторизация пользователей.
6. Получить избранных городов для пользователя.
7. Добавить город в избранное.
8. Удалить город из избранного.

Общее:
1. Приложение запускается в Docker-контейнере.
2. Написаны Unit-тесты для слоя репозитория и слоя бизнес-логики.
3. Добавлено логирование.
4. Добавлена Swagger-документация.
5. Описан makefile.

## Установка и запуск

1. Склонировать репозиторий.
```
git clone https://github.com/leerycorsair/weather-app.git
```

2. В папке docker/local приведен пример .env файла, который необходимо заполнить.
   
| **Флаг** | **Использование** | **Значение по умолчанию** | **Описание** |
|---|---|---|---|
| -s | -s | false | Включить получение данных из внешнего API. |
| -f | -f filename.txt | nil | Название файла, содержащего названия городов, которые будут загружены в сервис. |
| -u | -u 1m | 1m | Интервал обновления данных о погоде. |
| -p | -p | false | Включить параллельное получение данных. |

3. При первом запуске сервиса нужно обязательно создать текстовый файл с названиями городов, которые будут загружены в сервис (пример - cities.txt.example), а также задать соответствующие флаги.
   
4. Выполнить команду make run.

5. Swagger-документация будет доступна на localhost'е и соответствующем порте.

Видео демонстрация доступна по ссылке https://youtu.be/_S-nRVv5BWo.
   