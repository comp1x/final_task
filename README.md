## Реализация микросервисного приложения для приема заказов сотрудников офисов на обед.

### Сценарий использования:
1. Ресторан загружает меню на вторник в понедельник утром;
2. Сотрудники офисов получают это меню и создают заказы;
3. Вечером в понедельник запись закрывается;
4. Утром вторника ресторан собирает всю информация о заказах, готовит еду и отправляет службой доставки обеды в офисы.

### Что использовали.

1. GO
2. GRPC
3. ORM - GORM
4. БД - Postgres
5. Брокер - Apache Kafka

Текущий прогресс
- [x] CustomerService (CreateOffice, GetOfficeList, CreateOrder, GetActualMenu, GetUserList, CreateUser)
- [x] RestaurantService (CreateMenu, GetMenu, GetUpToDateOrderList, GetProductList, CreateProduct)
- [ ] StatisticsService (in future)

## Локальный запуск проекта


### Предварительная установка зависимостей.

```go mod tidy``` - установка зависимостей.



Перед началом работы нужно запустить docker-compose для запуска контейнеров, на котором находятся сервисы:

```docker-compose up```


### Применение миграций:
Выполняем:

```go run customer/pkg/migrate/migrate_up.go```

```go run restaurant/pkg/migrate/migrate_up.go```

```go run statistics/pkg/migrate/migrate_up.go```

### Запуск сервисов:

Запуск сервиса customer:

```go run customer/cmd/main/main.go```

Запуск сервиса restaurant:

```go run restaurant/cmd/main/main.go```

Запуск сервиса statistics:

```in future..```

### Специально для вас - [examples](../examples)

### В будущем:

1. Отрефакторить код)
2. Правильно обработать все ошибки
3. Дописать юнит тесты
4. Добавить logger (logrus) -> сохранять логи в стрим

