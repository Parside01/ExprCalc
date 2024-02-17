# Используем образ с Node.js current-alpine3.18
FROM node:13.12.0-alpine as build

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем файлы package.json и package-lock.json для установки зависимостей
COPY package*.json ./
COPY package.json ./
COPY package-lock.json ./

# Устанавливаем зависимости
RUN npm install

# Копируем все файлы из текущего каталога внутрь контейнера
COPY . .

# Запускаем сборку React приложения
RUN npm run build

# Этап production build
FROM nginx:1.21.3-alpine

# Копируем собранные файлы React приложения в Nginx
COPY --from=build /app/build /usr/share/nginx/html

# Открываем порт 80
EXPOSE 80

# Команда для запуска Nginx в режиме daemon
CMD ["nginx", "-g", "daemon off;"]