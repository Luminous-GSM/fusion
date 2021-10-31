package com.luminous.fusion.model.domain.server;

import com.fasterxml.jackson.annotation.JsonProperty;

public enum NodeStatus {
    @JsonProperty("running") RUNNING,
    @JsonProperty("pending") PENDING,
    @JsonProperty("terminated") TERMINATED,
    @JsonProperty("inactive") inactive,
}
