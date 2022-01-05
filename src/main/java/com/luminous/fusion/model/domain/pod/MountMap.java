package com.luminous.fusion.model.domain.pod;

import lombok.*;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
@ToString
public class MountMap {
    private String source;
    private String destination;
}
