FROM amazoncorretto:11-alpine-jdk

RUN mkdir /fusion
COPY fusion.jar /fusion/app.jar

EXPOSE 7878

#RUN apt-get -yqq update && apt-get -yqq install docker.io
VOLUME /var/run/docker.sock
VOLUME /fusion/plugins

ENTRYPOINT ["java","-jar","/fusion/app.jar"]