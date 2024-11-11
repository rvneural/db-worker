FROM debian:latest
LABEL maintainer="gafarov@realnoevremya.ru"
RUN apt-get update && apt-get upgrade
EXPOSE 7999
COPY . .
WORKDIR /build/linux
CMD [ "./dbworker" ]