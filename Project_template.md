## Изучите [README.md](.\README.md) файл и структуру проекта.

# Задание 1

1. Спроектируйте to be архитектуру КиноБездны, разделив всю систему на отдельные домены и организовав интеграционное взаимодействие и единую точку вызова сервисов.
Результат представьте в виде контейнерной диаграммы в нотации С4.
Добавьте ссылку на файл в этот шаблон
[C4 Container Diagram](c4.puml)

# Задание 2

### 1. Proxy
Команда КиноБездны уже выделила сервис метаданных о фильмах movies и вам необходимо реализовать бесшовный переход с применением паттерна Strangler Fig в части реализации прокси-сервиса (API Gateway), с помощью которого можно будет постепенно переключать траффик, используя фиче-флаг.


Реализуйте сервис на любом языке программирования в ./src/microservices/proxy.
Конфигурация для запуска сервиса через docker-compose уже добавлена
```yaml
  proxy-service:
    build:
      context: ./src/microservices/proxy
      dockerfile: Dockerfile
    container_name: cinemaabyss-proxy-service
    depends_on:
      - monolith
      - movies-service
      - events-service
    ports:
      - "8000:8000"
    environment:
      PORT: 8000
      MONOLITH_URL: http://monolith:8080
      #монолит
      MOVIES_SERVICE_URL: http://movies-service:8081 #сервис movies
      EVENTS_SERVICE_URL: http://events-service:8082 
      GRADUAL_MIGRATION: "true" # вкл/выкл простого фиче-флага
      MOVIES_MIGRATION_PERCENT: "50" # процент миграции
    networks:
      - cinemaabyss-network
```

#### Реализация прокси-сервиса

Прокси-сервис реализован на Go в `./src/microservices/proxy` и реализует паттерн **Strangler Fig** для постепенной миграции с монолита на микросервисы:

- HTTP-сервер слушает на порту 8000 и выступает единой точкой входа (API Gateway).
- Для endpoint `/api/movies` применяется **feature flag** через переменную `MOVIES_MIGRATION_PERCENT` (0–100). При каждом запросе генерируется случайное число 0–99: если оно меньше заданного процента — запрос проксируется на `movies-service`, иначе — на `monolith`.
- Переменная `GRADUAL_MIGRATION` позволяет полностью включить/выключить постепенную миграцию. Если выключена, весь трафик идёт на монолит.
- Запросы к `/api/events/*` маршрутизируются на `events-service`.
- Все остальные запросы (`/api/users`, `/api/payments`, `/api/subscriptions`) проксируются на `monolith`.
- Каждый проксированный запрос логируется с указанием, куда он был направлен.

- После реализации запустите postman тесты - они все должны быть зеленые (кроме events).
- Отправьте запросы к API Gateway:
   ```bash
   curl http://localhost:8000/api/movies
   ```
- Протестируйте постепенный переход, изменив переменную окружения MOVIES_MIGRATION_PERCENT в файле docker-compose.yml.


### 2. Kafka
 Вам как архитектуру нужно также проверить гипотезу насколько просто реализовать применение Kafka в данной архитектуре.

Для этого нужно сделать MVP сервис events, который будет при вызове API создавать и сам же читать сообщения в топике Kafka.

    - Разработайте сервис на любом языке программирования с consumer'ами и producer'ами.
    - Реализуйте простой API, при вызове которого будут создаваться события User/Payment/Movie и обрабатываться внутри сервиса с записью в лог
    - Добавьте в docker-compose новый сервис, kafka там уже есть

#### Реализация сервиса events

Сервис events реализован на Go в `./src/microservices/events` и интегрирован с Kafka:

- HTTP-сервер слушает на порту 8082 и реализует API согласно спецификации `api-specification.yaml`.
- **Эндпоинты:**
  - `GET /api/events/health` — health check
  - `POST /api/events/movie` — создание события фильма (MovieEvent)
  - `POST /api/events/user` — создание события пользователя (UserEvent)
  - `POST /api/events/payment` — создание события платежа (PaymentEvent)
- **Producer:** при вызове любого POST-эндпоинта сервис формирует сообщение с полями `id`, `type`, `timestamp`, `payload` и публикует его в Kafka-топик `cinema-events` (имя вынесено в env `KAFKA_TOPIC`).
- **Consumer:** при старте сервиса запускается фоновый consumer (group: `events-consumer-group`), который непрерывно читает сообщения из того же топика и записывает их в лог сервиса.
- Структура сообщения: `{ "id": "movie-1-viewed", "type": "movie", "timestamp": "...", "payload": {...} }`.
- В `docker-compose.yml` сервис добавлен с зависимостью от `kafka` и переменными `KAFKA_BROKERS` и `KAFKA_TOPIC`.
- Используется библиотека `segmentio/kafka-go`.

Необходимые тесты для проверки этого API вызываются при запуске npm run test:local из папки tests/postman 
Приложите скриншот тестов и скриншот состояния топиков Kafka из UI http://localhost:8090 

# Задание 3

Команда начала переезд в Kubernetes для лучшего масштабирования и повышения надежности. 
Вам, как архитектору осталось самое сложное:
 - реализовать CI/CD для сборки прокси сервиса
 - реализовать необходимые конфигурационные файлы для переключения трафика.


### CI/CD

#### Реализация CI/CD

Workflow `.github/workflows/docker-build-push.yml` доработан для сборки и деплоя всех четырёх сервисов:

1. **Триггеры:** push в ветки `main` и `cinema`, изменения в `src/**` или самом workflow, а также при создании релиза.
2. **Job `build-and-push`** собирает и пушит 4 Docker-образа в GitHub Container Registry (`ghcr.io`):
   - `ghcr.io/<repo>/monolith` из `./src/monolith`
   - `ghcr.io/<repo>/movies-service` из `./src/microservices/movies`
   - `ghcr.io/<repo>/proxy-service` из `./src/microservices/proxy`
   - `ghcr.io/<repo>/events-service` из `./src/microservices/events`
   - Теги: `sha-<short>`, `<branch>`, `latest`, семантические версии при релизе.
3. **Job `api-tests`** запускается после успешной сборки (`needs: build-and-push`):
   - Поднимает все сервисы через `docker compose up -d`
   - Ждёт готовности (120 сек)
   - Устанавливает Newman и запускает `npm run test:local`
   - При падении тестов workflow становится красным

 В папке .github/worflows доработайте деплой новых сервисов proxy и events в docker-build-push.yml , чтобы api-tests при сборке отрабатывали корректно при отправке коммита в ваш репозиторий.

Нужно доработать 
```yaml
on:
  push:
    branches: [ main ]
    paths:
      - 'src/**'
      - '.github/workflows/docker-build-push.yml'
  release:
    types: [published]
```
и добавить необходимые шаги в блок
```yaml
jobs:
  build-and-push:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
      - name: Checkout repository
        uses: actions/checkout@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Log in to the Container registry
        uses: docker/login-action@v2
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

```
Как только сборка отработает и в github registry появятся ваши образы, можно переходить к блоку настройки Kubernetes
Успешным результатом данного шага является "зеленая" сборка и "зеленые" тесты


#### Ресурсы Kubernetes для proxy и events

Созданы следующие манифесты:

- **`proxy-service.yaml`** — Deployment (1 реплика, порт 8000) + Service (ClusterIP). Env-переменные `MONOLITH_URL`, `MOVIES_SERVICE_URL`, `EVENTS_SERVICE_URL`, `GRADUAL_MIGRATION`, `MOVIES_MIGRATION_PERCENT` берутся из ConfigMap. Health-пробы на `/health`.
- **`events-service.yaml`** — Deployment (1 реплика, порт 8082) + Service (ClusterIP). Env-переменные `KAFKA_BROKERS` и `KAFKA_TOPIC` из ConfigMap. Health-пробы на `/api/events/health`.
- **`configmap.yaml`** — добавлены ключи `EVENTS_SERVICE_URL`, `KAFKA_BROKERS`, `KAFKA_TOPIC`. Ключ `MOVIES_MIGRATION_PERCENT` можно менять без пересборки образа — достаточно `kubectl rollout restart`.
- **`ingress.yaml`** — маршрутизация:
  - `cinemaabyss.example.com/api/events/*` → `events-service:8082` (напрямую, для Postman-тестов)
  - `cinemaabyss.example.com/*` → `proxy-service:8000` (API Gateway, Strangler Fig)
  - Запрос `https://cinemaabyss.example.com/api/movies` попадает в proxy, который по feature flag `MOVIES_MIGRATION_PERCENT` решает — отправить на `movies-service` или на `monolith`.

#### Инструкция запуска и тестирования (5 шагов)

1. Развернуть кластер: `kubectl apply -f src/kubernetes/` (namespace, secrets, configmap, postgres, kafka, сервисы, ingress).
2. Дождаться запуска подов: `kubectl -n cinemaabyss get pods` — все должны быть Running.
3. Включить ingress: `minikube addons enable ingress && minikube tunnel`.
4. Проверить: `curl https://cinemaabyss.example.com/api/movies` — должен вернуться список фильмов.
5. Запустить тесты: `cd tests/postman && npm run test:kubernetes`. Проверить логи events: `kubectl -n cinemaabyss logs -l app=events-service`.

### Proxy в Kubernetes

#### Шаг 1
Для деплоя в kubernetes необходимо залогиниться в docker registry Github'а.
1. Создайте Personal Access Token (PAT) https://github.com/settings/tokens . Создавайте class с правом read:packages
2. В src/kubernetes/*.yaml (event-service, monolith, movies-service и proxy-service)  отредактируйте путь до ваших образов 
```bash
 spec:
      containers:
      - name: events-service
        image: ghcr.io/ваш логин/имя репозитория/events-service:latest
```
3. Добавьте в секрет src/kubernetes/dockerconfigsecret.yaml в поле
```bash
 .dockerconfigjson: значение в base64 файла ~/.docker/config.json
```

4. Если в ~/.docker/config.json нет значения для аутентификации
```json
{
        "auths": {
                "ghcr.io": {
                       тут пусто
                }
        }
}
```
то выполните 

и добавьте

```json 
 "auth": "имя пользователя:токен в base64"
```

Чтобы получить значение в base64 можно выполнить команду
```bash
 echo -n ваш_логин:ваш_токен | base64
```

После заполнения config.json, также прогоните содержимое через base64

```bash
cat .docker/config.json | base64
```

и полученное значение добавляем в

```bash
 .dockerconfigjson: значение в base64 файла ~/.docker/config.json
```

#### Шаг 2

  Доработайте src/kubernetes/event-service.yaml и src/kubernetes/proxy-service.yaml

  - Необходимо создать Deployment и Service 
  - Доработайте ingress.yaml, чтобы можно было с помощью тестов проверить создание событий
  - Выполните дальшейшие шаги для поднятия кластера:

  1. Создайте namespace:
  ```bash
  kubectl apply -f src/kubernetes/namespace.yaml
  ```
  2. Создайте секреты и переменные
  ```bash
  kubectl apply -f src/kubernetes/configmap.yaml
  kubectl apply -f src/kubernetes/secret.yaml
  kubectl apply -f src/kubernetes/dockerconfigsecret.yaml
  kubectl apply -f src/kubernetes/postgres-init-configmap.yaml
  ```

  3. Разверните базу данных:
  ```bash
  kubectl apply -f src/kubernetes/postgres.yaml
  ```

  На этом этапе если вызвать команду
  ```bash
  kubectl -n cinemaabyss get pod
  ```
  Вы увидите

  NAME         READY   STATUS    
  postgres-0   1/1     Running   

  4. Разверните Kafka:
  ```bash
  kubectl apply -f src/kubernetes/kafka/kafka.yaml
  ```

  Проверьте, теперь должно быть запущено 3 пода, если что-то не так, то посмотрите логи
  ```bash
  kubectl -n cinemaabyss logs имя_пода (например - kafka-0)
  ```

  5. Разверните монолит:
  ```bash
  kubectl apply -f src/kubernetes/monolith.yaml
  ```
  6. Разверните микросервисы:
  ```bash
  kubectl apply -f src/kubernetes/movies-service.yaml
  kubectl apply -f src/kubernetes/events-service.yaml
  ```
  7. Разверните прокси-сервис:
  ```bash
  kubectl apply -f src/kubernetes/proxy-service.yaml
  ```

  После запуска и поднятия подов вывод команды 
  ```bash
  kubectl -n cinemaabyss get pod
  ```

  Будет наподобие такого

```bash
  NAME                              READY   STATUS    

  events-service-7587c6dfd5-6whzx   1/1     Running  

  kafka-0                           1/1     Running   

  monolith-8476598495-wmtmw         1/1     Running  

  movies-service-6d5697c584-4qfqs   1/1     Running  

  postgres-0                        1/1     Running  

  proxy-service-577d6c549b-6qfcv    1/1     Running  

  zookeeper-0                       1/1     Running 
```

  8. Добавим ingress

  - добавьте аддон
  ```bash
  minikube addons enable ingress
  ```
  ```bash
  kubectl apply -f src/kubernetes/ingress.yaml
  ```
  9. Добавьте в /etc/hosts
  127.0.0.1 cinemaabyss.example.com

  10. Вызовите
  ```bash
  minikube tunnel
  ```
  11. Вызовите https://cinemaabyss.example.com/api/movies
  Вы должны увидеть вывод списка фильмов
  Можно поэкспериментировать со значением   MOVIES_MIGRATION_PERCENT в src/kubernetes/configmap.yaml и убедится, что вызовы movies уходят полностью в новый сервис

  12. Запустите тесты из папки tests/postman
  ```bash
   npm run test:kubernetes
  ```
  Часть тестов с health-чек упадет, но создание событий отработает.
  Откройте логи event-service и сделайте скриншот обработки событий

#### Шаг 3
Добавьте сюда скриншота вывода при вызове https://cinemaabyss.example.com/api/movies и  скриншот вывода event-service после вызова тестов.


# Задание 4

### Описание Helm-чарта `cinema`

Helm-чарт расположен в `src/kubernetes/helm/` и разворачивает все компоненты «Кинобездны» одной командой.

**Создаваемые ресурсы:**
- **Deployment + Service** для `proxy-service` (Strangler Fig API Gateway, порт 8000)
- **Deployment + Service** для `events-service` (Kafka producer/consumer, порт 8082)
- **Deployment + Service** для `monolith` и `movies-service`
- **StatefulSet** для PostgreSQL, Kafka и Zookeeper
- **ConfigMap** `cinemaabyss-config` — конфигурация для всех сервисов
- **Ingress** — маршрутизация через `cinemaabyss.example.com`
- **Secrets** — Docker registry credentials и DB password

**Ключевые параметры в `values.yaml`:**

| Параметр | Описание | Значение по умолчанию |
|---|---|---|
| `proxyService.image.repository` | Образ proxy-сервиса | `ghcr.io/shatim-ops/architecture-cinemaabyss/proxy-service` |
| `proxyService.image.tag` | Тег образа proxy | `latest` |
| `proxyService.image.pullPolicy` | Pull policy proxy | `Always` |
| `proxyService.env.MOVIES_MIGRATION_PERCENT` | Процент миграции (feature flag) | `50` |
| `eventsService.image.repository` | Образ events-сервиса | `ghcr.io/shatim-ops/architecture-cinemaabyss/events-service` |
| `eventsService.image.tag` | Тег образа events | `latest` |
| `eventsService.image.pullPolicy` | Pull policy events | `Always` |
| `eventsService.env.KAFKA_BROKERS` | Адрес Kafka-брокера | `kafka:9092` |
| `eventsService.env.KAFKA_TOPIC` | Топик Kafka | `cinema-events` |
| `config.moviesMigrationPercent` | MOVIES_MIGRATION_PERCENT в ConfigMap | `50` |
| `config.kafkaBrokers` | KAFKA_BROKERS в ConfigMap | `kafka:9092` |
| `config.kafkaTopic` | KAFKA_TOPIC в ConfigMap | `cinema-events` |
| `ingress.enabled` | Включить Ingress | `true` |
| `ingress.hosts[0].host` | Хост Ingress | `cinemaabyss.example.com` |
| `proxyService.service.type` | Тип Service proxy | `ClusterIP` |
| `eventsService.service.type` | Тип Service events | `ClusterIP` |

**Как установить/обновить чарт:**

```bash
# Установка (или обновление)
helm upgrade --install cinema ./src/kubernetes/helm -n cinemaabyss --create-namespace

# Проверка подов
kubectl get pods,svc,ingress -n cinemaabyss

# Изменение MOVIES_MIGRATION_PERCENT без пересборки
helm upgrade cinema ./src/kubernetes/helm -n cinemaabyss \
  --set config.moviesMigrationPercent=100

# Или через values-файл
helm upgrade cinema ./src/kubernetes/helm -n cinemaabyss -f custom-values.yaml
```

**Проверка работоспособности после установки:**

1. Убедиться, что все поды Running: `kubectl get pods -n cinemaabyss`
2. Включить ingress и tunnel: `minikube addons enable ingress && minikube tunnel`
3. Проверить маршрут: `curl https://cinemaabyss.example.com/api/movies` — должен вернуть список фильмов
4. Изменить процент миграции: `helm upgrade cinema ./src/kubernetes/helm -n cinemaabyss --set config.moviesMigrationPercent=100` — весь трафик `api/movies` пойдёт на `movies-service`
5. Запустить тесты: `cd tests/postman && npm run test:kubernetes`, проверить логи events: `kubectl -n cinemaabyss logs -l app=events-service`

---

Для простоты дальнейшего обновления и развертывания вам как архитектуру необходимо так же реализовать helm-чарты для прокси-сервиса и проверить работу 

Для этого:
1. Перейдите в директорию helm и отредактируйте файл values.yaml

```yaml
# Proxy service configuration
proxyService:
  enabled: true
  image:
    repository: ghcr.io/db-exp/cinemaabysstest/proxy-service
    tag: latest
    pullPolicy: Always
  replicas: 1
  resources:
    limits:
      cpu: 300m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi
  service:
    port: 80
    targetPort: 8000
    type: ClusterIP
```

- Вместо ghcr.io/db-exp/cinemaabysstest/proxy-service напишите свой путь до образа для всех сервисов
- для imagePullSecret проставьте свое значение (скопируйте из конфигурации kubernetes)
  ```yaml
  imagePullSecrets:
      dockerconfigjson: ewoJImF1dGhzIjogewoJCSJnaGNyLmlvIjogewoJCQkiYXV0aCI6ICJaR0l0Wlhod09tZG9jRjl2UTJocVZIa3dhMWhKVDIxWmFVZHJOV2hRUW10aFVXbFZSbTVaTjJRMFNYUjRZMWM9IgoJCX0KCX0sCgkiY3JlZHNTdG9yZSI6ICJkZXNrdG9wIiwKCSJjdXJyZW50Q29udGV4dCI6ICJkZXNrdG9wLWxpbnV4IiwKCSJwbHVnaW5zIjogewoJCSIteC1jbGktaGludHMiOiB7CgkJCSJlbmFibGVkIjogInRydWUiCgkJfQoJfSwKCSJmZWF0dXJlcyI6IHsKCQkiaG9va3MiOiAidHJ1ZSIKCX0KfQ==
  ```

2. В папке ./templates/services заполните шаблоны для proxy-service.yaml и events-service.yaml (опирайтесь на свою kubernetes конфигурацию - смысл helm'а сделать шаблоны для быстрого обновления и установки)

```yaml
template:
    metadata:
      labels:
        app: proxy-service
    spec:
      containers:
       Тут ваша конфигурация
```

3. Проверьте установку
Сначала удалим установку руками

```bash
kubectl delete all --all -n cinemaabyss
kubectl delete  namespace cinemaabyss
```
Запустите 
```bash
helm install cinemaabyss .\src\kubernetes\helm --namespace cinemaabyss --create-namespace
```
Если в процессе будет ошибка
```code
[2025-04-08 21:43:38,780] ERROR Fatal error during KafkaServer startup. Prepare to shutdown (kafka.server.KafkaServer)
kafka.common.InconsistentClusterIdException: The Cluster ID OkOjGPrdRimp8nkFohYkCw doesn't match stored clusterId Some(sbkcoiSiQV2h_mQpwy05zQ) in meta.properties. The broker is trying to join the wrong cluster. Configured zookeeper.connect may be wrong.
```

Проверьте развертывание:
```bash
kubectl get pods -n cinemaabyss
minikube tunnel
```

Потом вызовите 
https://cinemaabyss.example.com/api/movies
и приложите скриншот развертывания helm и вывода https://cinemaabyss.example.com/api/movies

## Удаляем все

```bash
kubectl delete all --all -n cinemaabyss
kubectl delete namespace cinemaabyss
```
