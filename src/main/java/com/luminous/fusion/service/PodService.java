package com.luminous.fusion.service;

import com.github.dockerjava.api.DockerClient;
import com.github.dockerjava.api.command.CreateContainerResponse;
import com.github.dockerjava.api.command.ListContainersCmd;
import com.github.dockerjava.api.command.PullImageResultCallback;
import com.github.dockerjava.api.model.*;
import com.luminous.fusion.model.request.pod.PodCreateRequest;
import com.luminous.fusion.model.request.pod.PodRemoveRequest;
import com.luminous.fusion.model.request.pod.PodStartRequest;
import com.luminous.fusion.model.request.pod.PodStopRequest;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.stream.Collectors;

@Service
@AllArgsConstructor
@Slf4j
public class PodService {

    private final DockerClient dockerClient;

    public String getServerVersion() {
        return this.dockerClient.infoCmd().exec().getServerVersion();
    }

    public void initialiseServer() {

    }

    public void createPod(PodCreateRequest podCreateRequest) throws InterruptedException {

        this.dockerClient
                .pullImageCmd(podCreateRequest.getPodDescription().getImage())
                .withTag(podCreateRequest.getPodDescription().getTag())
                .exec(new PullImageResultCallback())
                .awaitCompletion();

        log.info("Docker - Pull Image Completed");

        HostConfig hostConfig = new HostConfig()
                .withPortBindings(
                        podCreateRequest
                                .getPodDescription()
                                .getPortMaps()
                                .stream()
                                .map(portMap ->
                                        new PortBinding(
                                                Ports.Binding.bindPort(portMap.getBinding()),
                                                new ExposedPort(portMap.getExposed(), InternetProtocol.parse(portMap.getProtocol()))
                                        )
                                )
                                .collect(Collectors.toList())
                )
                .withMounts(
                        podCreateRequest
                                .getPodDescription()
                                .getMountMaps()
                                .stream()
                                .map(volumeMap ->
                                        new Mount()
                                                .withSource(volumeMap.getSource())
                                                .withTarget(volumeMap.getDestination())
                                                .withType(MountType.BIND)
                                )
                                .collect(Collectors.toList())
                )
                .withRestartPolicy(RestartPolicy.unlessStoppedRestart());


        CreateContainerResponse createContainerResponse = this.dockerClient
                .createContainerCmd(podCreateRequest.getPodDescription().getImage())
                .withName(podCreateRequest.getPodDescription().getName())
                .withHostConfig(hostConfig)
                .withExposedPorts(
                        podCreateRequest
                                .getPodDescription()
                                .getPortMaps()
                                .stream()
                                .map(portMap ->
                                        new ExposedPort(portMap.getExposed(), InternetProtocol.parse(portMap.getProtocol()))
                                )
                                .collect(Collectors.toList())
                )
                .withEnv(
                        podCreateRequest
                                .getPodDescription()
                                .getEnvironmentMaps()
                                .stream()
                                .map(environmentMap ->
                                        String.format("%s=%s",environmentMap.getName(), environmentMap.getValue())
                                )
                                .collect(Collectors.toList())
                )
                .exec();

        log.info("Docker - Create Container Completed");

        this.dockerClient.startContainerCmd(createContainerResponse.getId()).exec();

        log.info("Docker - Start Container Completed");
    }

    public void stopPod(PodStopRequest podStopRequest) {
        this.dockerClient.stopContainerCmd(podStopRequest.getContainerId()).exec();
    }

    public void startPod(PodStartRequest podStartRequest) {
        this.dockerClient.startContainerCmd(podStartRequest.getContainerId()).exec();
    }

    public List<Container> listContainers(boolean showALl, Integer exitedFilter) {
        ListContainersCmd listContainersCmd = this.dockerClient.listContainersCmd();

        listContainersCmd.withShowAll(exitedFilter != null || showALl);
        if (exitedFilter != null) {
            listContainersCmd.withExitedFilter(exitedFilter);
        }

        return listContainersCmd.exec();
    }

    public void removePod(PodRemoveRequest podRemoveRequest) {
        this.dockerClient
                .removeContainerCmd(podRemoveRequest.getContainerId())
                .withForce(podRemoveRequest.isForceRemove())
                .exec();
    }

}
