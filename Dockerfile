FROM umputun/baseimage:buildgo-v1.9.2 as build
ADD . /build
WORKDIR  /build
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /target/pt

FROM umputun/baseimage:app-v1.9.2
COPY --from=build /target/pt /srv/pt
RUN chown -R app:app /srv
WORKDIR /srv
ENTRYPOINT ["/srv/pt"]
