FROM amazoncorretto:11-alpine-jdk

RUN mkdir /fusion
COPY fusion.jar /fusion/app.jar

EXPOSE 7878

RUN apk update && apk add --no-cache docker-cli
VOLUME /var/run/docker.sock
VOLUME /fusion/plugins

ENTRYPOINT ["java -XX:+UseSerialGC -Xss512k -XX:MaxRAM=100m","-jar","/fusion/app.jar"]