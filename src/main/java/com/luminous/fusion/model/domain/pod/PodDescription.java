package com.luminous.fusion.model.domain.pod;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.Set;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class PodDescription {
    private String id;
    private String name;
    private String image;
    private String tag;
    private Set<PortMap> portMaps;
    private Set<EnvironmentMap> environmentMaps;
    private Set<MountMap> mountMaps;
}
