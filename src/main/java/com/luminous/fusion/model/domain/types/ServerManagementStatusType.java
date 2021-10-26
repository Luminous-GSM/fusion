package com.luminous.fusion.model.domain.types;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.AllArgsConstructor;
import lombok.Getter;

@Getter
public enum ServerManagementStatusType {
    @JsonProperty("uninitialized") UNINITIALIZED,
    @JsonProperty("initialized") INITIALIZED
}
