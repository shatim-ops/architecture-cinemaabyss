# Кинобездна — проектная работа спринта 2

Проект демонстрирует переход от монолитной системы к микросервисной архитектуре, настройку инфраструктуры (Kubernetes, Service Mesh, Kafka), а также CI/CD и Helm‑чарты для сервисов прокси и событий. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)

## Структура репозитория

- `docs/` — документация проекта.  
  - `cinema-architecture-c4-container.puml` — C4‑диаграмма уровня контейнеров To‑Be архитектуры «Кинобездны». [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- `src/microservices/` — микросервисы проекта.  
  - `proxy/` — прокси‑сервис (API‑Gateway / BFF) для маршрутизации запросов от клиентов к монолиту и новому сервису фильмов. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
  - `events/` — сервис событий (producer/consumer) для интеграции с Kafka. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
  - другие сервисы и заглушки, предоставленные Практикумом. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- `src/kubernetes/` — Kubernetes‑манифесты (до Helm).  
  - `proxy-service.yaml` — Deployment и Service для прокси‑сервиса. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
  - `event-service.yaml` — Deployment и Service для сервиса событий. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
  - `configmap.yaml` — ConfigMap с конфигурацией (в т.ч. `MOVIES_MIGRATION_PERCENT`). [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
  - `ingress.yaml` — Ingress‑ресурс для маршрутизации `https://cinemaabyss.example.com`. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- `helm/cinema/` — Helm‑чарт для развёртывания proxy и events в Kubernetes.  
  - `Chart.yaml` — метаданные чарта. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
  - `values.yaml` — настройки образов, портов, ingress и конфигурации. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
  - `templates/` — шаблоны Deployment, Service, ConfigMap и Ingress. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- `.github/workflows/docker-build-push.yml` — Pipeline as Code для сборки и публикации образов, запуска API‑тестов. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- `Project_template.md` — шаблон отчёта по проекту с описанием решений по заданиям 1–4. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)

## Архитектура (To‑Be)

- Клиентские приложения (Web, Mobile, Smart TV) обращаются к единой точке входа — прокси‑сервису (BFF/API‑Gateway). [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- Прокси‑сервис маршрутизирует запросы к монолиту или новым микросервисам (например, movies) с использованием паттерна Strangler Fig и feature flag. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- Сервис событий (`events`) публикует доменные события (просмотры, платежи, подписки) в Kafka и читает их обратно, логируя обработку. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- Микросервисы развернуты в Kubernetes, доступ организован через Service Mesh и Ingress, конфигурация вынесена в ConfigMap и Helm‑values. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)

Подробная схема указана в `docs/cinema-architecture-c4-container.puml`. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)

## Задание 1: C4‑диаграмма архитектуры

- Спроектирована To‑Be архитектура «Кинобездны» с декомпозицией на микросервисы (пользователи, платежи, подписки, каталог, скидки, события и интеграции). [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- Диаграмма уровня контейнеров оформлена в формате PlantUML C4 и размещена в `docs/cinema-architecture-c4-container.puml`. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)

## Задание 2: Proxy и Events с Docker/Kafka

### Proxy‑сервис

- Расположение: `src/microservices/proxy/`. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- Функции:  
  - поднимает HTTP‑сервер с endpoint’ом `GET /api/movies`;  
  - читает переменную окружения `MOVIES_MIGRATION_PERCENT` (0–100);  
  - для каждого запроса генерирует случайное число и по значению percent решает, отправлять запрос в новый `movies`‑сервис или в старый монолитный endpoint (паттерн Strangler Fig);  
  - логирует маршрутизацию запроса. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)

### Events‑сервис

- Расположение: `src/microservices/events/`. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- Функции:  
  - реализует HTTP‑API по спецификации проекта (создание событий пользователей/платежей/просмотров);  
  - при приёме запросов публикует сообщения в Kafka‑топик;  
  - внутри сервиса работает consumer, который читает сообщения из Kafka и пишет информацию в лог. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)

### Инфраструктура

- Docker‑контейнеры описаны в `docker-compose.yml`, включая Kafka и UI для Kafka. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- Postman‑тесты, приложенные в проекте, используются для проверки работы proxy/events; ожидается успешное прохождение всех тестов, кроме специально помеченных health‑checks. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)

## Задание 3: CI/CD и Kubernetes‑манифесты

### CI/CD (GitHub Actions)

- Файл workflow: `.github/workflows/docker-build-push.yml`. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- Основные шаги:  
  - сборка Docker‑образов для `proxy` и `events` из `src/microservices/proxy` и `src/microservices/events`;  
  - публикация образов в контейнерный Registry (на основе настроек репозитория);  
  - запуск API‑тестов (Postman/newman) после сборки; пайплайн падает при неуспешных тестах. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)

### Kubernetes‑манифесты

- `src/kubernetes/proxy-service.yaml`  
  - Deployment для прокси‑сервиса: образ из Registry, переменные окружения, probes;  
  - Service (обычно ClusterIP) для внутреннего доступа к proxy. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- `src/kubernetes/event-service.yaml`  
  - Deployment для сервиса событий: образ, конфигурация подключения к Kafka, probes;  
  - Service для доступа к events внутри кластера. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- `src/kubernetes/configmap.yaml`  
  - ConfigMap с конфигурацией, включая `MOVIES_MIGRATION_PERCENT` и другие параметры;  
  - подключается к proxy через `env`/`envFrom`. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- `src/kubernetes/ingress.yaml`  
  - Ingress‑правило для маршрута `https://cinemaabyss.example.com/api/movies` на proxy‑service;  
  - дополнительные пути для API events‑сервиса (используются Postman‑тестами). [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)

## Задание 4: Helm‑чарт

### Структура чарта

- Каталог чарта: `helm/cinema/`. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
- Основные файлы:  
  - `Chart.yaml` — описание чарта (имя `cinema`, версия `0.1.0`).  
  - `values.yaml` — параметры образов, портов, ingress и конфигурации; управляет `MOVIES_MIGRATION_PERCENT`. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
  - `templates/proxy-deployment.yaml` / `proxy-service.yaml` — шаблоны Deployment и Service для proxy. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
  - `templates/events-deployment.yaml` / `events-service.yaml` — шаблоны Deployment и Service для events. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
  - `templates/configmap.yaml` — шаблон ConfigMap с конфигами proxy. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)
  - `templates/ingress.yaml` — шаблон Ingress‑ресурса. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)

### Использование Helm

- Рендеринг шаблонов:  
  ```bash
  helm template cinema ./helm/cinema
  ```  
- Установка/обновление:  
  ```bash
  helm upgrade --install cinema ./helm/cinema -n cinema --create-namespace
  ```  
- Изменение процента миграции трафика:  
  - через `values.yaml` (`proxy.env.MOVIES_MIGRATION_PERCENT`);  
  - либо флагом:  
    ```bash
    helm upgrade cinema ./helm/cinema \
      -n cinema \
      --set proxy.env.MOVIES_MIGRATION_PERCENT=100
    ``` [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)  

## Локальный запуск и тестирование

### Docker Compose

1. Установить Docker и Docker Compose.  
2. Запустить окружение:  
   ```bash
   docker-compose up --build
   ```  
3. Проверить ручной запрос:  
   ```bash
   curl http://localhost:<port>/api/movies
   ```  
4. Запустить Postman‑тесты из набора, предоставленного Практикумом, и убедиться, что тесты для proxy и events проходят успешно. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)

### Kubernetes + Helm

1. Поднять локальный Kubernetes‑кластер (minikube/kind/k3d по инструкции Практикума).  
2. Установить чарт:  
   ```bash
   helm upgrade --install cinema ./helm/cinema -n cinema --create-namespace
   ```  
3. Проверить ресурсы:  
   ```bash
   kubectl get pods,svc,ingress -n cinema
   ```  
4. Настроить доступ к `https://cinemaabyss.example.com` (hosts/порт‑форвардинг) и проверить:  
   ```bash
   curl https://cinemaabyss.example.com/api/movies
   ``` [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)  

## Отчёт по проекту

Подробные ответы по заданиям 1–4, а также ссылки на диаграмму, скриншоты Kafka/Postman и описание пайплайна CI/CD находятся в файле `Project_template.md` в корне репозитория. [practicum.yandex](https://practicum.yandex.ru/learn/software-architect/courses/ee3aa7a8-3f51-463f-85ad-b62c2c4feaf4/sprints/810032/topics/7242c699-f94c-4d1c-987b-11ab12b094a1/lessons/0511caa5-ea42-423c-b36c-d6a596e0701c/)

