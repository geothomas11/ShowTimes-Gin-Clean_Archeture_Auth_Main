FROM golang:1.24-alpine AS build-stage
WORKDIR /app
COPY ./ /app
RUN mkdir -p /app/build
RUN go mod download
RUN go build -v -o /app/build/api ./cmd/api


FROM gcr.io/distroless/static-debian11
COPY --from=build-stage /app/build/api /api
COPY --from=build-stage /app/pkg/templates /app/templates
COPY --from=build-stage /app/.env /
EXPOSE 7000
CMD ["/api"] 