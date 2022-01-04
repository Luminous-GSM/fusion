# The base_image should hold a reference to the image created by ./gradlew bootBuildImage
ARG base_image
FROM ${base_image}

RUN apt-get -yqq update && apt-get -yqq install docker.io
VOLUME /var/run/docker.sock

RUN echo "Hello Custom Dockerfile"

ENTRYPOINT /cnb/process/web