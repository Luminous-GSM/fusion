package com.luminous.fusion.controller;

import com.luminous.fusion.service.PodHealthService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/health")
public class HealthController {

    private final PodHealthService podHealthService;

    @Autowired
    public HealthController(PodHealthService podHealthService) {
        this.podHealthService = podHealthService;
    }

    @GetMapping("/status")
    public String getHealthForFusion() {
        return "Up";
    }

    @GetMapping("/plugin/{id}")
    public String getHealthForPodByIdViaPlugin(@PathVariable String id) {
        return podHealthService.getHealthForPodById(id);
    }


}
