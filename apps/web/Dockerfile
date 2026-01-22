FROM node:24-alpine AS builder

WORKDIR /app

COPY package.json pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install --frozen-lockfile

COPY . .
RUN pnpm build

# Production: Nginx serving static files on port 8080
FROM bitnami/nginx:latest

COPY --from=builder /app/dist /app
EXPOSE 8080
CMD ["nginx", "-g", "daemon off;"]
