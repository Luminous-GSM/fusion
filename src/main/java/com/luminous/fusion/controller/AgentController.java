package com.luminous.fusion.controller;

import com.luminous.fusion.service.AgentService;
import com.luminous.fusion.service.PodService;
import lombok.AllArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Map;

@RestController
@RequestMapping("/agent")
@AllArgsConstructor
public class AgentController {

    private final AgentService agentService;
    private final PodService podService;

    @GetMapping("/ping")
    public ResponseEntity<String> pingPong() {
        return ResponseEntity.ok("Pong");
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
