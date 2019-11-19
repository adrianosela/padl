FROM alpine:latest
RUN apk add --update bash curl && rm -rf /var/cache/apk/*
COPY . .
EXPOSE 80
CMD ["./padl"]
