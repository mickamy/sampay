ARG NODE_VERSION=22.13.0

FROM node:${NODE_VERSION}-alpine AS development-dependencies-env
COPY . /src
WORKDIR /src
RUN npm ci

FROM node:${NODE_VERSION}-alpine AS production-dependencies-env
COPY ./package.json package-lock.json /src/
WORKDIR /src
RUN npm ci --omit=dev

FROM node:${NODE_VERSION}-alpine AS build-env
COPY . /src/
COPY --from=development-dependencies-env /src/node_modules /src/node_modules
WORKDIR /src
RUN npm run build

FROM node:${NODE_VERSION}-alpine
COPY ./package.json package-lock.json /src/
COPY --from=production-dependencies-env /src/node_modules /src/node_modules
COPY --from=build-env /src/build /src/build
WORKDIR /src
CMD ["npm", "run", "start"]
