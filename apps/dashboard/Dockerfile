FROM node:20-alpine as builder

WORKDIR /app

COPY package*.json ./

RUN corepack enable

RUN pnpm install --frozen-lockfile

COPY . .

RUN pnpm run build

FROM oven/bun:1.3-alpine

WORKDIR /app

COPY --from=builder /app/dist .

EXPOSE 3001

CMD ["bun", "serve", "-p", "3001", "."]
