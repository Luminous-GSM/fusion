package com.luminous.fusion.controller;

import com.github.dockerjava.api.model.Container;
import com.luminous.fusion.model.request.pod.PodCreateRequest;
import com.luminous.fusion.model.request.pod.PodRemoveRequest;
import com.luminous.fusion.model.request.pod.PodStartRequest;
import com.luminous.fusion.model.request.pod.PodStopRequest;
import com.luminous.fusion.service.PodService;
import lombok.AllArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Optional;

@RestController
@RequestMapping("/pods")
@AllArgsConstructor
public class PodsController {

    private final PodService podService;

    @PostMapping("/")
    public ResponseEntity<Object> createPod(@RequestBody PodCreateRequest podCreateRequest) {
        try {
            this.podService.createPod(podCreateRequest);

            return ResponseEntity.ok("Ok");
        } catch (InterruptedException e) {
            return ResponseEntity.internalServerError().body(e.getMessage());
        }
    }

    @GetMapping("/")
    public ResponseEntity<List<Container>> getAllPods(
            @RequestParam(defaultValue = "false", required = false) boolean showAll,
            @RequestParam(required = false) Integer exitedFilter) {
        return ResponseEntity.ok(
                this.podService.listContainers(showAll, exitedFilter)
        );
    }

    @DeleteMapping("/")
    public ResponseEntity<String> removePod(@RequestBody PodRemoveRequest podRemoveRequest) {
        this.podService.removePod(podRemoveRequest);
        return ResponseEntity.ok(String.format("%s removed successfully", podRemoveRequest.getContainerId()));
    }

    @PostMapping("/stop")
    public ResponseEntity<String> stopPod(@RequestBody PodStopRequest podStopRequest) {
        this.podService.stopPod(podStopRequest);
        return ResponseEntity.ok(String.format("%s stopped successfully", podStopRequest.getContainerId()));
    }

    @PostMapping("/start")
    public ResponseEntity<String> startPod(@RequestBody PodStartRequest podStartRequest) {
        this.podService.startPod(podStartRequest);
        return ResponseEntity.ok(String.format("%s started successfully", podStartRequest.getContainerId()));
    }

}
