ARG DOCKER_BUILD_IMAGE=golang:1.13

FROM ${DOCKER_BUILD_IMAGE} AS build
WORKDIR /app/
COPY . /app/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build .

# Final Image
FROM iron/base
LABEL name="Mattermost Birthday Bot" \
  maintainer="ctadeu@gmail.com"

WORKDIR /app

COPY --from=build /app/gcalendar /app

ENTRYPOINT /app/gcalendar
