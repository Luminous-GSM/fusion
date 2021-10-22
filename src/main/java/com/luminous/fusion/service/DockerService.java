package com.luminous.fusion.service;

import com.github.dockerjava.api.DockerClient;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;

@Service
@AllArgsConstructor
public class DockerService {

    private final DockerClient dockerClient;

    public String getServerVersion() {
        return this.dockerClient.infoCmd().exec().getServerVersion();
    }

}
