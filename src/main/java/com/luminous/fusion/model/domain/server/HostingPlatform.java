package com.luminous.fusion.model.domain.server;

import com.fasterxml.jackson.annotation.JsonProperty;

public enum HostingPlatform {
    @JsonProperty("local") LOCAL,
    @JsonProperty("docker") DOCKER;


}
