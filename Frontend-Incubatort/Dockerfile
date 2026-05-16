# Stage 1: build Vue app
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Stage 2: serve dengan Nginx
FROM nginx:alpine

# Copy hasil build ke direktori Nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# Config Nginx untuk Vue Router (history mode) + proxy ke backend
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
