FROM oven/bun as build-stage

WORKDIR /app

COPY package.json ./
COPY bun.lock ./
RUN bun install
COPY . .

ARG VUE_APP_API_URL
ENV VITE_API_URL=${VUE_APP_API_URL}

RUN bun run build

FROM nginx:stable-alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
COPY ./nginx-default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]