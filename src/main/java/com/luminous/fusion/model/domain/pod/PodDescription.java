package com.luminous.fusion.model.domain.pod;

import lombok.*;

import java.util.Set;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
@ToString
public class PodDescription {
    private String id;
    private String name;
    private String image;
    private String tag;
    private Set<PortMap> portMaps;
    private Set<EnvironmentMap> environmentMaps;
    private Set<MountMap> mountMaps;
    private String command;
    private String manifestFileUsed;
}
