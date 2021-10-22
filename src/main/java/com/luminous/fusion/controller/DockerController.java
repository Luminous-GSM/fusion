package com.luminous.fusion.controller;

import com.luminous.fusion.service.DockerService;
import lombok.AllArgsConstructor;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/docker")
@AllArgsConstructor
public class DockerController {

    private final DockerService dockerService;

    @GetMapping("/server/version")
    public String getServerVersion() {
        return this.dockerService.getServerVersion();
    }
}
