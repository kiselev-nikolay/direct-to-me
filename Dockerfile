FROM node AS frontend
WORKDIR /tmp/src
COPY ./assets/src/package.json .
COPY ./assets/src/package-lock.json .
RUN npm i
COPY ./assets/ /tmp/
RUN npm run build

FROM golang:1.16 AS backend_builder
WORKDIR /build
COPY ./pkg/ ./pkg
COPY ./cmd/ ./cmd
COPY ./go.mod ./
COPY ./go.sum ./
RUN go build ./cmd/server/main.go

FROM golang:1.16 AS backend
WORKDIR /srv/app
COPY --from=frontend /tmp/public /srv/app/assets/public
COPY --from=backend_builder /build/main /srv/app
CMD /srv/app/main
