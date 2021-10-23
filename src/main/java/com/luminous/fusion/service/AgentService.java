package com.luminous.fusion.service;

import com.luminous.fusion.configuration.LuminousPropertiesConfiguration;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@AllArgsConstructor
public class AgentService {

    private final LuminousPropertiesConfiguration luminousPropertiesConfiguration;

    public String getVersion() {
        return this.luminousPropertiesConfiguration.getVersion();
    }

    public void initializeServer() {

    }

}
