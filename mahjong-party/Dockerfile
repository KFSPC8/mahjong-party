# base image
FROM node:16-alpine AS build

WORKDIR /app

# copy only package.json and lock for caching
COPY package.json yarn.lock ./

RUN yarn install

# copy rest and build
COPY . .
RUN yarn build

# production image
FROM nginx:alpine

COPY --from=build /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]