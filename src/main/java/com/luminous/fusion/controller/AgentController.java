package com.luminous.fusion.controller;

import com.luminous.fusion.configuration.LuminousPropertiesConfiguration;
import com.luminous.fusion.model.domain.server.NodeDescription;
import com.luminous.fusion.model.domain.server.NodeStatus;
import com.luminous.fusion.model.response.agent.ContainerDto;
import com.luminous.fusion.model.response.agent.DashboardResponse;
import com.luminous.fusion.model.response.agent.ImageDto;
import com.luminous.fusion.model.response.agent.SystemLoadResponse;
import com.luminous.fusion.service.AgentService;
import com.luminous.fusion.service.PodService;
import com.sun.management.OperatingSystemMXBean;
import lombok.AllArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.io.IOException;
import java.lang.management.ManagementFactory;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

@RestController
@RequestMapping("/agent")
@AllArgsConstructor
public class AgentController {

    private final AgentService agentService;
    private final PodService podService;

    private final LuminousPropertiesConfiguration luminousPropertiesConfiguration;

    @GetMapping("/ping")
    public ResponseEntity<NodeDescription> pingPong() {
        return ResponseEntity.ok(
                new NodeDescription(
                        luminousPropertiesConfiguration.getNode().getUniqueId(),
                        luminousPropertiesConfiguration.getNode().getName(),
                        luminousPropertiesConfiguration.getNode().getDescription(),
                        luminousPropertiesConfiguration.getVersion(),
                        luminousPropertiesConfiguration.getPlatform(),
                        NodeStatus.RUNNING,
                        this.podService.getTotalActivePods()
                )
        );
    }

    @GetMapping("/system-load")
    public ResponseEntity<SystemLoadResponse> getSystemLoad() {
        OperatingSystemMXBean osBean = ManagementFactory.getPlatformMXBean(OperatingSystemMXBean.class);

        return ResponseEntity.ok(
                new SystemLoadResponse(
                        Math.round(osBean.getSystemCpuLoad() * 100),
                        Math.round( ( (double) (osBean.getTotalPhysicalMemorySize() - osBean.getFreePhysicalMemorySize() ) / (double) osBean.getTotalPhysicalMemorySize() ) * 100),
                        Math.round(0)
                )
        );
    }

    @GetMapping("/dashboard")
    public ResponseEntity<DashboardResponse> getAgentDashboard() {
        return ResponseEntity.ok(
                new DashboardResponse(
                        new NodeDescription(
                                luminousPropertiesConfiguration.getNode().getUniqueId(),
                                luminousPropertiesConfiguration.getNode().getName(),
                                luminousPropertiesConfiguration.getNode().getDescription(),
                                luminousPropertiesConfiguration.getVersion(),
                                luminousPropertiesConfiguration.getPlatform(),
                                NodeStatus.RUNNING,
                                this.podService.getTotalActivePods()
                        ),
                        this.podService.listContainers(List.of()).stream().map(ContainerDto::new).collect(Collectors.toList()),
                        this.podService.getImages().stream().map(ImageDto::new).collect(Collectors.toList())
                )
        );
    }



}
