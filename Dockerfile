# The base_image should hold a reference to the image created by ./gradlew bootBuildImage
ARG base_image
FROM ${base_image}

VOLUME /var/run/docker.sock

ENTRYPOINT /cnb/process/web