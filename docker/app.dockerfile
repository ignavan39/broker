FROM node:14.18-alpine as dev

WORKDIR /app/
COPY ./yarn.lock ./package.json ./

RUN yarn install --frozen-lockfile

COPY . .

FROM node:14.18-alpine as builder

WORKDIR /app/

COPY --from=dev /app/ /app/

RUN yarn build


FROM node:14.18-alpine

WORKDIR /app/

COPY --from=builder /app/package.json ./
COPY --from=builder /app/yarn.lock ./

RUN yarn install --production --frozen-lockfile

COPY --from=builder /app/config ./config
COPY --from=builder /app/dist ./dist