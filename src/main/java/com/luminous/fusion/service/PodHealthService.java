package com.luminous.fusion.service;

import interfaces.game.Health;
import org.pf4j.PluginManager;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class PodHealthService {

    private PluginManager pluginManager;

    public String getHealthForPodById(String id) {
        // TODO Get plugin for name for pod.

        List podHealths = this.pluginManager.getExtensions(id);

        return "Shap";
    }

}
