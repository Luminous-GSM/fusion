package com.luminous.fusion.model.domain.pod;


import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class EnvironmentMap {
    private String name;
    private String value;
}
