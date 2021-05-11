# build ui stage
FROM node:lts-alpine as ui-build-stage
WORKDIR /app
COPY ui/package*.json ./
RUN npm install
COPY ui .
RUN npm run build && ls -lash dist

FROM golang:buster as api-build-stage
WORKDIR /api

COPY api/ .
RUN GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags='-w -s -extldflags "-static"' -a \
        -o plooral-api ./cmd/plooral \
        && chmod a+x plooral-api

FROM nginx:latest
COPY --from=ui-build-stage /app/dist /app/ui
COPY --from=api-build-stage /api/plooral-api /app/plooral-api
COPY ui/nginx.conf /etc/nginx/nginx.conf
EXPOSE 8079
CMD nginx -g 'daemon off;' & /app/plooral-api