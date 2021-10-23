package com.luminous.fusion.controller;

import com.luminous.fusion.service.PluginService;
import lombok.AllArgsConstructor;
import models.health.HealthResult;
import models.statistics.StatisticsResult;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/plugin")
@AllArgsConstructor
public class PluginController {

    private final PluginService pluginService;

    @GetMapping("/reload")
    public ResponseEntity<Object> reloadPlugins() {
        this.pluginService.reloadPlugins();

        return ResponseEntity.ok("Ok");
    }

    @GetMapping("/unload")
    public ResponseEntity<Object> unloadPlugins() {
        this.pluginService.unloadPlugins();

        return ResponseEntity.ok("Ok");
    }

    @GetMapping("/load")
    public ResponseEntity<Object> loadPlugins() {
        this.pluginService.loadPlugins();

        return ResponseEntity.ok("Ok");
    }

    @GetMapping("/{id}/HEALTH")
    public ResponseEntity<HealthResult> getHealthViaPlugin(@PathVariable String id) {
        HealthResult healthResult = this.pluginService.getHealthViaPlugin(id);

        return ResponseEntity.ok(healthResult);
    }

    @GetMapping("/{id}/STATISTICS")
    public ResponseEntity<StatisticsResult> getStatisticsViaPlugin(@PathVariable String id) {
        StatisticsResult statisticsResult = this.pluginService.getStatisticsViaPlugin(id);

        return ResponseEntity.ok(statisticsResult);
    }

}
