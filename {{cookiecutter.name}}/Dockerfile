# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /build/{{ cookiecutter.name }} ./cmd/{{ cookiecutter.name }}

# Runtime stage - distroless for security
FROM gcr.io/distroless/static-debian11:nonroot
COPY --from=builder /build/{{ cookiecutter.name }} /
USER 65534
ENTRYPOINT ["/{{ cookiecutter.name }}"]

LABEL org.opencontainers.image.title="{{ cookiecutter.name }}"
LABEL org.opencontainers.image.source="https://github.com/{{ cookiecutter.repo_owner }}/{{ cookiecutter.repo_name }}"
