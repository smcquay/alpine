FROM alpine:3.5
RUN apk add --no-cache curl
RUN apk add --no-cache bind-tools
RUN apk add --no-cache jq
ADD bin/servedir /usr/local/bin/
ADD bin/cs /usr/local/bin/
ENTRYPOINT ["sleep", "3600"]
