package com.luminous.fusion.model.request.pod;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.lang.NonNull;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public abstract class ContainerId {
    @NonNull
    private String containerId;
}
