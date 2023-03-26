# Этап 1: Компиляция двоичного файла в контейнеризованном окружении Golang
# Использовать образ "golang:1.19-alpine" для компиляции. Компиляция происходит в контейнере "build"
FROM golang:1.19-alpine as builder
# Назначить рабочим каталог с исходным кодом
WORKDIR /go/src/go_ping_kube
#COPY go.mod .
#COPY go.sum .
#RUN go mod download
# Скопировать исходные файлы из локального контекста в рабочую директорию образа
COPY ./ ./
# Собрать двоичный файл
RUN CGO_ENABLED=0 GOOS=linux go build -v -o ./bin/go_ping_kube ./cmd/ping_server

# Этап 2: - Сборка образа
# Использовать образ "alpine"
FROM alpine
# Назначить рабочим каталог /app
WORKDIR /app
RUN apk --no-cache add ca-certificates
# Скопировать двоичный файл из контейнера build
COPY --from=builder /go/src/go_ping_kube/bin/go_ping_kube /app/go_ping_kube
# Скопировать конфиг
COPY ./configs /app/configs
# Команда, которая должна быть выполнена при запуске контейнера
CMD ["/app/go_ping_kube"]