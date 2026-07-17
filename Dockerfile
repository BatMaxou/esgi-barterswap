FROM golang:1.26.3

ARG GID=1000
ARG UID=1000

RUN (groupadd -g $GID go || true) && \
    useradd -d /home/go -g $GID -u $UID -m go

USER go
WORKDIR /home/go

EXPOSE 8000

CMD ["go", "run", "./src"]
