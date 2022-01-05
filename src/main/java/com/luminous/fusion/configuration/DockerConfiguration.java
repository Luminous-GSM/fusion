package com.luminous.fusion.configuration;

import com.github.dockerjava.api.DockerClient;
import com.github.dockerjava.core.DefaultDockerClientConfig;
import com.github.dockerjava.core.DockerClientConfig;
import com.github.dockerjava.core.DockerClientImpl;
import com.github.dockerjava.httpclient5.ApacheDockerHttpClient;
import com.github.dockerjava.transport.DockerHttpClient;
import com.luminous.fusion.model.domain.server.HostingPlatform;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.DependsOn;

@Configuration
@AllArgsConstructor
@Slf4j
public class DockerConfiguration {

    private final LuminousPropertiesConfiguration luminousPropertiesConfiguration;

    private String getHostnameBasedOnHostingPlatform() {
        log.info(String.valueOf(luminousPropertiesConfiguration));

        String host = "";
        switch (luminousPropertiesConfiguration.getPlatform()) {
            case LOCAL: host = "tcp://127.0.0.1:2375";
            break;
            case DOCKER: host = "tcp://host.docker.internal:2375";
            break;
            case AWS: host = String.format("tcp://%s:2375", luminousPropertiesConfiguration.getNode().getHostname());
            break;
        }
        log.info(host);
        return host;
    }

    @Bean
    public DockerClientConfig dockerClientConfig() {
        return DefaultDockerClientConfig.createDefaultConfigBuilder()
                .withDockerHost(getHostnameBasedOnHostingPlatform())
                .build();
    }

    @Bean
    @DependsOn("dockerClientConfig")
    public DockerHttpClient dockerHttpClient(DockerClientConfig dockerClientConfig) {
        return new ApacheDockerHttpClient.Builder()
                .dockerHost(dockerClientConfig.getDockerHost())
                .sslConfig(dockerClientConfig.getSSLConfig())
                .build();
    }

    @Bean
    @DependsOn("dockerHttpClient")
    public DockerClient dockerClient(DockerClientConfig dockerClientConfig, DockerHttpClient dockerHttpClient) {
        return DockerClientImpl.getInstance(dockerClientConfig, dockerHttpClient);
    }

}
