# Build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json ./
RUN npm install --legacy-peer-deps
COPY . .
# Vite attend index.html a la racine
RUN cp public/index.html . 2>/dev/null || true
RUN npm run build

# Production stage (nginx)
FROM nginx:alpine AS prod
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

# Dev stage
FROM node:20-alpine AS dev
WORKDIR /app
COPY package.json ./
RUN npm install --legacy-peer-deps
COPY . .
RUN cp public/index.html . 2>/dev/null || true
EXPOSE 3000
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
