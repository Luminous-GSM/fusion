package com.luminous.fusion.model.response.agent;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.github.dockerjava.api.model.ContainerPort;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class ContainerPortDto {

    private String ip;
    private Integer privatePort;
    private Integer publicPort;
    private String type;

    public ContainerPortDto(ContainerPort containerPort) {
        this.ip = containerPort.getIp();
        this.privatePort = containerPort.getPrivatePort();
        this.publicPort = containerPort.getPublicPort();
        this.type = containerPort.getType();
    }
}
