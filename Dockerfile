FROM alpine:3.5
RUN apk add --no-cache curl
RUN apk add --no-cache bind-tools
RUN apk add --no-cache jq
ADD bin/servedir /usr/local/bin/
ADD bin/cs /usr/local/bin/
ENTRYPOINT ["/bin/sh", "-c", "while true; do sleep 3600; done"]
