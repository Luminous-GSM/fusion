package com.luminous.fusion.service;

import interfaces.game.Health;
import interfaces.game.Statistics;
import lombok.AllArgsConstructor;
import models.health.HealthResult;
import models.statistics.StatisticsResult;
import org.pf4j.PluginManager;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
@AllArgsConstructor
public class PluginService {
    private PluginManager pluginManager;

    public HealthResult getHealthViaPlugin(String id) {

        List<Health> plugins = this.pluginManager.getExtensions(Health.class, id);

        return plugins.get(0).health();
    }

    public StatisticsResult getStatisticsViaPlugin(String id) {

        List<Statistics> plugins = this.pluginManager.getExtensions(Statistics.class, id);

        return plugins.get(0).statistics();
    }

    public void reloadPlugins() {
        this.pluginManager.unloadPlugins();
        this.pluginManager.loadPlugins();
    }

    public void loadPlugins() {
        this.pluginManager.loadPlugins();
        // TODO Return a list of plugins in their specific state
    }

    public void unloadPlugins() {
        this.pluginManager.unloadPlugins();
    }

    public void unloadPlugin(String pluginId) {
        this.pluginManager.unloadPlugin(pluginId);
    }
}
