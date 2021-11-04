package com.luminous.fusion.model.response.agent;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class SystemLoadResponse {
    private double cpuLoad;
    private double ramLoad;
    private double hddUsage;
}
