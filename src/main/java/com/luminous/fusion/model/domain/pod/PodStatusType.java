package com.luminous.fusion.model.domain.pod;

import lombok.AllArgsConstructor;
import lombok.Getter;

@AllArgsConstructor
@Getter
public enum PodStatusType {
    CREATED("created"),
    RESTARTING("restarting"),
    RUNNING("running"),
    REMOVING("removing"),
    PAUSED("paused"),
    EXITED("exited"),
    DEAD("dead");

    private String value;
}
