package com.luminous.fusion.service;

import com.github.dockerjava.api.DockerClient;
import com.github.dockerjava.api.command.CreateContainerResponse;
import com.github.dockerjava.api.command.ListContainersCmd;
import com.github.dockerjava.api.command.LogContainerCmd;
import com.github.dockerjava.api.command.PullImageResultCallback;
import com.github.dockerjava.api.model.*;
import com.luminous.fusion.callbacks.FusionLogContainerCallback;
import com.luminous.fusion.model.request.pod.PodCreateRequest;
import com.luminous.fusion.model.request.pod.PodRemoveRequest;
import com.luminous.fusion.model.request.pod.PodStartRequest;
import com.luminous.fusion.model.request.pod.PodStopRequest;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;

@Service
@AllArgsConstructor
@Slf4j
public class PodService {

    private final String LABEL_IS_FUSION_POD = "is-fusion-pod";
    private final int WAIT_TIMEOUT = 5;

    private final DockerClient dockerClient;

    public String getServerVersion() {
        return this.dockerClient.infoCmd().exec().getServerVersion();
    }

    public void initialiseServer() {

    }

    public String createPod(PodCreateRequest podCreateRequest) throws InterruptedException {
        log.info("Service-PodService | createPod | Start");

        log.info("Docker - Pull {} | Starting", podCreateRequest.getPodDescription().getImage());

        this.dockerClient
                .pullImageCmd(podCreateRequest.getPodDescription().getImage())
                .withTag(podCreateRequest.getPodDescription().getTag())
                .exec(new PullImageResultCallback())
                .awaitCompletion();

        log.info("Docker - Pull {} | Completed", podCreateRequest.getPodDescription().getImage());

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

        log.info("Docker - Create Container {} | Starting", podCreateRequest.getPodDescription().getName());

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
                                        String.format("%s=%s", environmentMap.getName(), environmentMap.getValue())
                                )
                                .collect(Collectors.toList())
                )
                .withLabels(
                        Map.of(
                                "manifest-file-used", podCreateRequest.getPodDescription().getManifestFileUsed(),
                                LABEL_IS_FUSION_POD, "true"
                        )
                )
                .exec();

        log.info("Docker - Create Container {} | Completed", podCreateRequest.getPodDescription().getName());

        log.info("Service-PodService | createPod | End");
        return createContainerResponse.getId();

    }

    public Container getSimpleContainerViaId(String containerId) {
        return this.dockerClient.listContainersCmd().withShowAll(true).withIdFilter(List.of(containerId)).exec().get(0);
    }

    public void stopPod(PodStopRequest podStopRequest) {
        this.dockerClient.stopContainerCmd(podStopRequest.getContainerId()).withTimeout(WAIT_TIMEOUT).exec();
    }

    public void startPod(PodStartRequest podStartRequest) {
        this.dockerClient.restartContainerCmd(podStartRequest.getContainerId()).withtTimeout(WAIT_TIMEOUT).exec();
    }

    public List<Container> listContainers(boolean showALl, Integer exitedFilter) {
        ListContainersCmd listContainersCmd = this.dockerClient.listContainersCmd();

        listContainersCmd.withShowAll(exitedFilter != null || showALl);
        if (exitedFilter != null) {
            listContainersCmd.withExitedFilter(exitedFilter);
        }

        return listContainersCmd.exec();
    }

    public List<Container> listContainers(List<String> status) {

        ListContainersCmd listContainersCmd = this.dockerClient
                .listContainersCmd()
                .withShowAll(true)
                .withLabelFilter(
                        Map.of(LABEL_IS_FUSION_POD, "true")
                );

        if (!status.isEmpty()) {
            listContainersCmd.withStatusFilter(status);
        }

        return listContainersCmd.exec();
    }

    public List<Image> getImages() {
        return this.dockerClient
                .listImagesCmd()
                .withShowAll(true)
                .exec();
    }

    public void removePod(PodRemoveRequest podRemoveRequest) {
        this.dockerClient
                .removeContainerCmd(podRemoveRequest.getContainerId())
                .withForce(podRemoveRequest.isForceRemove())
                .exec();
    }

    public String getContainerLogs(String containerId) throws InterruptedException {

        FusionLogContainerCallback fusionLogContainerCallback = new FusionLogContainerCallback(true);

        this.dockerClient
                .logContainerCmd(containerId)
                .withStdOut(true)
                .withStdErr(true)
                .withFollowStream(true)
                .withTailAll()
                .withTimestamps(true)
                .exec(fusionLogContainerCallback);

        fusionLogContainerCallback.awaitCompletion(3, TimeUnit.SECONDS);

        return fusionLogContainerCallback.toString();
    }

    public int getTotalActivePods() {
        return this.dockerClient.listContainersCmd().exec().size();
    }

}
