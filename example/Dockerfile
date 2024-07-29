FROM golang AS builder
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./example

FROM alpine:latest
ARG USER=default
ENV HOME /home/$USER
RUN adduser -D $USER
COPY --from=builder /server $HOME/server
RUN chown -R $USER:$USER $HOME
USER $USER
WORKDIR $HOME
ENTRYPOINT [ "/home/default/server" ]