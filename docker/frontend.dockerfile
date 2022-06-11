FROM node:17-alpine as builder

WORKDIR /app

COPY ./frontend/package.json ./
COPY ./frontend/yarn.lock ./

RUN yarn install --frozen-lockfile

COPY ./frontend/ .

RUN yarn build

FROM nginx:1.19.9

COPY ./docker/nginx/frontend.nginx.conf /etc/nginx/nginx.conf
RUN rm -rf /usr/share/nginx/html/*

COPY --from=builder /app/build /usr/share/nginx/html
