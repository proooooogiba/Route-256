
### Примеры grp-curl-запросов

### Комната
##### Создание:
`grpcurl -plaintext -d '{"name":"Example Room","cost":100.00}' localhost:50051 hotel.HotelService/CreateRoom`

##### Получение:
`grpcurl -plaintext -d '{"id":1}' localhost:50051 hotel.HotelService/GetRoomWithAllReservations`

##### Обновление:
`grpcurl -plaintext -d '{"id":1, "name":"Updated Example Room", "cost":100.0}' localhost:50051 hotel.HotelService/UpdateRoom`

##### Удаление:
`grpcurl -plaintext -d '{"id":1}' localhost:50051 hotel.HotelService/DeleteRoom`


### Бронирование
##### Создание:
`grpcurl -plaintext -d '{"start_date":"2023-09-15", "end_date":"2023-11-15", "room_id":1}' localhost:50051 hotel.HotelService/CreateReservation`

##### Получение:
`grpcurl -plaintext -d '{"id":1}' localhost:50051 hotel.HotelService/GetReservation`

##### Обновление:
`grpcurl -plaintext -d '{"id":1, "start_date":"2023-09-15", "end_date":"2023-11-15", "room_id":1}' localhost:50051 hotel.HotelService/UpdateReservation`

##### Удаление:
`grpcurl -plaintext -d '{"id":1}' localhost:50051 hotel.HotelService/DeleteReservation`
