package com.luminous.fusion.configuration;

import com.github.dockerjava.api.DockerClient;
import com.github.dockerjava.core.DefaultDockerClientConfig;
import com.github.dockerjava.core.DockerClientConfig;
import com.github.dockerjava.core.DockerClientImpl;
import com.github.dockerjava.httpclient5.ApacheDockerHttpClient;
import com.github.dockerjava.transport.DockerHttpClient;
import com.luminous.fusion.model.domain.server.HostingPlatform;
import lombok.AllArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.DependsOn;

@Configuration
@AllArgsConstructor
public class DockerConfiguration {

    private final LuminousPropertiesConfiguration luminousPropertiesConfiguration;

    private String getHostnameBasedOnHostingPlatform() {
        switch (luminousPropertiesConfiguration.getPlatform()) {
            case LOCAL: return "tcp://127.0.0.1:2375";
            case DOCKER: return "tcp://host.docker.internal:2375";
            case AWS: return luminousPropertiesConfiguration.getNode().getHostname();
        }
        return "";
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
