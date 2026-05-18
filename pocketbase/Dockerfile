FROM alpine:latest
RUN apk add --no-cache ca-certificates unzip wget
RUN wget https://github.com/pocketbase/pocketbase/releases/download/v0.22.21/pocketbase_0.22.21_linux_amd64.zip &&     unzip pocketbase_0.22.21_linux_amd64.zip -d /usr/local/bin/ &&     rm pocketbase_0.22.21_linux_amd64.zip
EXPOSE 8090
CMD ["pocketbase", "serve", "--http=0.0.0.0:8090", "--dir=/pb_data", "--publicDir=/pb_public"]
