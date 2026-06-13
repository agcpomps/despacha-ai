FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

RUN npm run build

FROM node:22-alpine

WORKDIR /app

ENV NODE_ENV=production

# instala apenas as dependências de runtime (ex: valibot) que o adapter-node
# não empacota no bundle e que o servidor precisa de resolver em runtime
COPY package.json package-lock.json ./
RUN npm ci --omit=dev

COPY --from=builder /app/build ./build

EXPOSE 3000

CMD ["node", "build"]
