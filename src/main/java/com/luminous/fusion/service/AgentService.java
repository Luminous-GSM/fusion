package com.luminous.fusion.service;

import com.luminous.fusion.configuration.LuminousPropertiesConfiguration;
import interfaces.BasePlugin;
import interfaces.game.Health;
import lombok.AllArgsConstructor;
import org.pf4j.PluginManager;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@AllArgsConstructor
public class AgentService {

    private final LuminousPropertiesConfiguration luminousPropertiesConfiguration;

    public String getVersion() {
        return this.luminousPropertiesConfiguration.getVersion();
    }

}
