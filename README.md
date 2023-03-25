# OpenBlocks "Ограничитель запросов"

### Инициатива OpenBlocks

Инициатива OpenBlocks &mdash; это проект с открытым исходным кодом, целью которого
является предоставить открытые и масштабируемые решения уровня предприятия.

### Описание

Сервис "Ограничитель запросов" предоставляет возможность ограничить количество запросов,
отправляемых к какому-то ресурсу в единицу времени.

Этот вариант сервиса основывается на оригинальной [реализации на Java](https://github.com/IgorIvkin/openblocks-ratelimiter),
но написан на языке программирования **Go**. Реализация на Go несколько быстрее стартует, потребляет меньше памяти и в целом
является более легковесной.

Сервис спроектирован для применения в составе информационных систем, к нему должны обращаться
другие сервисы, чтобы определить, не ограничен ли доступ к какому-то ресурсу.

Механизм распределенного ограничения доступа
основывается на алгоритме **Token bucket**, по которому раз в единицу времени на ресурс выделяется 
некоторое заданное количество токенов доступа, а каждый запрос использует один токен. Если все токены
исчерпались, это означает, что доступ к ресурсу ограничен до последующего выпуска новых токенов.

### Основная конфигурация

Обратите внимание на следующие секции в конфигурации сервиса.

```yaml
limiters:
  basic:
    limit: 10
    unit: MINUTES
  test:
    limit: 5
    unit: SECONDS
```

В этом блоке заданы два рейт-лимитера, один с названием `basic` и с ограничением доступа 10 раз в минуту,
другой &mdash; с названием `test` и с ограничением доступа 5 раз в секунду.

Вы можете задать сколько угодно рейт-лимитеров, добавляя новые секции. Обратите внимание, каждая секция
должна иметь уникальное имя, которое будет являться ключом для доступа к этому рейт-лимитеру.


### API

GET /api/v1/rate-limits/{limiterName}

Проверяет возможность обращения к ресурсу по ключу рейт-лимитера &mdash; `limiterName`. Например, для рейт-лимитера
"basic" значение параметра будет равно "basic".

Возвращает `true`, если доступ возможен, и `false`, если доступ временно ограничен.



## Полезные ссылки
* [OpenBlocks "Пользователи"](https://github.com/IgorIvkin/openblocks-users)
* [OpenBlocks "Роли"](https://github.com/IgorIvkin/openblocks-roles)
* [OpenBlocks "Команды"](https://github.com/IgorIvkin/openblocks-teams)
* [Сервис "Ограничитель запросов" на Java](https://github.com/IgorIvkin/openblocks-ratelimiter)
* [Документация по Keycloak](https://www.keycloak.org/documentation)