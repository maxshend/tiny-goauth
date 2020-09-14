FROM migrate/migrate

RUN apk update && apk add --no-cache make

COPY /migrations /migrations
COPY /Makefile /Makefile

ENTRYPOINT ["/bin/sh", "-c"]
