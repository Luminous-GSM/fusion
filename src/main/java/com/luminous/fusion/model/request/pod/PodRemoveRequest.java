package com.luminous.fusion.model.request.pod;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class PodRemoveRequest extends ContainerId {
    private boolean forceRemove;
}
