package com.luminous.fusion.model.response.agent;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.github.dockerjava.api.command.ListContainersCmd;
import com.github.dockerjava.api.model.*;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.Arrays;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class ContainerDto {

    private String command;
    private Long created;
    private String id;
    private String image;
    private String imageId;
    private String[] names;
    public List<ContainerPortDto> ports;
    public Map<String, String> labels;
    private String status;
    private String state;
    private Long sizeRw;
    private Long sizeRootFs;
    private ContainerHostConfig hostConfig;
    private ContainerNetworkSettings networkSettings;
    private List<ContainerMount> mounts;

    public ContainerDto(Container container) {
        this.command = container.getCommand();
        this.created = container.getCreated();
        this.id = container.getId();
        this.image = container.getImage();
        this.imageId = container.getImageId();
        this.names = container.getNames();
        this.ports = Arrays.stream(container.getPorts()).map(ContainerPortDto::new).collect(Collectors.toList());
        this.labels = container.getLabels();
        this.status = container.getStatus();
        this.state = container.getState();
        this.sizeRw = container.getSizeRw();
        this.sizeRootFs = container.getSizeRootFs();
        this.hostConfig = container.getHostConfig();
        this.networkSettings = container.getNetworkSettings();
        this.mounts = container.getMounts();
    }


}
