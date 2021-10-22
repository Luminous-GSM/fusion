package com.luminous.fusion.controller;

import com.luminous.fusion.service.AgentService;
import lombok.AllArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/agent")
@AllArgsConstructor
public class AgentController {

    private final AgentService agentService;

    @GetMapping("/version")
    public String getAgentVersion() {
        return this.agentService.getVersion();
    }


}
