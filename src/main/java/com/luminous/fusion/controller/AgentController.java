package com.luminous.fusion.controller;

import com.luminous.fusion.configuration.LuminousPropertiesConfiguration;
import com.luminous.fusion.model.domain.server.NodeDescription;
import com.luminous.fusion.model.domain.server.NodeStatus;
import com.luminous.fusion.model.response.management.ManagementPingResult;
import com.luminous.fusion.service.AgentService;
import com.luminous.fusion.service.PodService;
import lombok.AllArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.io.IOException;
import java.util.Map;

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
                        NodeStatus.RUNNING)
        );
    }

    @GetMapping("/version")
    public ResponseEntity<Map<String, String>> getAgentVersion() {
        return ResponseEntity.ok(
                Map.of(
                        "agent", this.agentService.getVersion(),
                        "docker", this.podService.getServerVersion()
                )
        );
    }

    @GetMapping("/initialise")
    public ResponseEntity<Object> initialiseServer() {

        this.agentService.initializeServer();
        this.podService.initialiseServer();

        return ResponseEntity.ok("Ok");
    }



}
