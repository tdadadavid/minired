FROM golang AS build
ADD . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go install .

FROM alpine
RUN apk --no-cache add ca-certificates
COPY --from=build /go/bin/minired /minired
ENTRYPOINT ["/minired"]