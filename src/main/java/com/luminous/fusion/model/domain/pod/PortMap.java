package com.luminous.fusion.model.domain.pod;

import lombok.*;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
@ToString
public class PortMap {
    private Integer exposed; // Port that is internet facing
    private Integer binding; // Port that the container uses
    private String protocol;
}
