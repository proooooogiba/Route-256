### Даталогическая модель базы данных

![даталогическая модель бд](https://i.postimg.cc/wBQnZfTt/image.png)

### Примеры curl-запросов

При GET-запросе комнаты отдается она и все её бронирования, при DELETE удалется она и все её бронирования.
### Комната
##### Создание:
```curl -X POST localhost:9000/room -d '{"name":"Lux", "cost":1000.0}' -i```

##### Получение:
`curl -X GET localhost:9000/room/1 -i`

##### Обновление:
`curl -X PUT localhost:9000/room -d '{"id":1, "name":"NedoLux", "cost":91.0}' -i`

##### Удаление:
`curl -X DELETE localhost:9000/room/1 -i`


### Бронирование
##### Создание:
`curl -X POST localhost:9000/reservation -d '{"start_date":"2023-09-15", "end_date":"2023-11-15", "room_id":1}' -i`

##### Получение:
`curl -X GET localhost:9000/reservation/1 -i`

##### Обновление:
`curl -X PUT localhost:9000/reservation -d '{"id":1, "start_date":"2006-01-02", "end_date":"2006-01-07", "room_id":1}' -i`

##### Удаление:
`curl -X DELETE localhost:9000/reservation/1 -i`
