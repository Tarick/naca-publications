FROM golang:1.15.6-alpine as build
RUN apk --no-cache add tzdata curl && \
    curl --silent -o /tmp/swaggerui.tgz https://codeload.github.com/swagger-api/swagger-ui/tar.gz/v3.26.1 && \
    tar -C /tmp  -xvzf /tmp/swaggerui.tgz swagger-ui-3.26.1/dist --strip-components=1 && mv /tmp/dist /tmp/swaggerui &&\
        sed -i -e 's|https://petstore.swagger.io/v2/swagger.json|http://localhost:8080/doc/swagger.json|g' /tmp/swaggerui/index.html
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN GOOS=linux GOARCH=amd64 go mod download
ENV PACKAGE=github.com/Tarick/naca-publications
ARG BUILD_VERSION
ADD . .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags 'netgo static_build' -ldflags "-extldflags=-static -w -s -X ${PACKAGE}/pkg/version.Version=$BUILD_VERSION -X ${PACKAGE}/pkg/version.BuildTime=$(date --utc  +%FT%TZ)" -o ./build/publications-api ./cmd/publications-api

FROM scratch as final
WORKDIR /
ENV TZ=UTC
COPY --from=build /tmp/swaggerui /swaggerui
COPY --from=build /app/internal/docs/swagger.json /swaggerui
COPY --from=build /app/build/publications-api /
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
ENTRYPOINT ["./publications-api"]
CMD [ "--config", "./config.yaml" ]
