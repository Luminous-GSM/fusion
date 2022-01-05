package com.luminous.fusion.model.domain.pod;


import lombok.*;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
@ToString
public class EnvironmentMap {
    private String name;
    private String value;
}
