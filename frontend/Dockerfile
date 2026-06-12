FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# PUBLIC_* vars são embebidas no build do cliente
ARG PUBLIC_API_BASE_URL
ENV PUBLIC_API_BASE_URL=$PUBLIC_API_BASE_URL
ARG PUBLIC_ADMIN_WHATSAPP
ENV PUBLIC_ADMIN_WHATSAPP=$PUBLIC_ADMIN_WHATSAPP

RUN npm run build

FROM node:22-alpine

WORKDIR /app

ENV NODE_ENV=production

COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./

EXPOSE 3000

CMD ["node", "build"]
