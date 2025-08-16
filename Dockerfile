FROM golang:1.24.5 AS build

WORKDIR /app
RUN go env -w GOMODCACHE=/root/.cache/go-build

# Download Go modules
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
# RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY *.go ./

# Build
# RUN CGO_ENABLED=0 GOOS=linux go build -o homeservice
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 go build -o homeservice

FROM gcr.io/distroless/static-debian12 AS runner

COPY --from=build /app/homeservice /

EXPOSE 8080

CMD ["/homeservice"]