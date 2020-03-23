FROM golang:1.14.0 as base_builder

ENV TGT_CA_BUNDLE_PATH=/usr/local/share/ca-certificates/tgt-ca-bundle.crt
RUN mkdir -p /usr/local/etc /usr/local/share/ca-certificates && \
	curl -o $TGT_CA_BUNDLE_PATH http://browserconfig.target.com/tgt-certs/tgt-ca-bundle.crt && \
	cp $TGT_CA_BUNDLE_PATH /etc/ssl/certs/ca-certificates.crt && \
	update-ca-certificates 
    
# now create a minimal image that will build the application.
# FROM alpine:3.7
FROM docker.target.com/tap/alpine-adoptopenjdk:11
EXPOSE 8080
# COPY --from=base_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENV CLOUD_ENVIRONMENT=prod
WORKDIR /root/

COPY ./vault-env-linux-amd64-0.0.2 ./vault-env
COPY ./runtime.yml .
COPY ./secret.tpl.yml /tap/config/secret.tpl.yml
COPY ./vault-secret.yml /tap/secret/vault-secret.yml

RUN chmod u+x ./vault-env
ENTRYPOINT ["./vault-env render"]