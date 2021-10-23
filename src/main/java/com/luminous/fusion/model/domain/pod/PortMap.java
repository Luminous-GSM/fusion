package com.luminous.fusion.model.domain.pod;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class PortMap {
    private Integer exposed; // Port that is internet facing
    private Integer binding; // Port that the container uses
    private String protocol;
}
