package com.luminous.fusion.service;

import interfaces.BasePlugin;
import lombok.AllArgsConstructor;
import org.pf4j.PluginManager;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@AllArgsConstructor
public class PluginService {
    private PluginManager pluginManager;

    public String getHealthForPodById(String id) {

        List<BasePlugin> plugins = this.pluginManager.getExtensions(BasePlugin.class, id);

        return plugins.get(0).health().getStatus();
    }
}
