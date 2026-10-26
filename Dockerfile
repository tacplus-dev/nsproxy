FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY api api
COPY cmd cmd
COPY internal internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags=-s -o /nsproxy ./cmd/nsproxy

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /nsproxy /nsproxy
USER 65532:65532
ENTRYPOINT ["/nsproxy"]
