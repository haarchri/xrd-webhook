FROM golang:1.20.3-alpine AS build

WORKDIR /src
COPY /go.mod /go.sum ./
RUN go mod download
COPY ./ ./

RUN CGO_ENABLED=false go build --installsuffix 'static' -o xrd-webhook .

FROM alpine AS final
COPY --from=build /src/xrd-webhook /xrd-webhook
