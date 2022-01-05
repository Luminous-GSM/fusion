package com.luminous.fusion.controller;

import com.github.dockerjava.api.model.Container;
import com.luminous.fusion.model.request.pod.PodCreateRequest;
import com.luminous.fusion.model.request.pod.PodRemoveRequest;
import com.luminous.fusion.model.request.pod.PodStartRequest;
import com.luminous.fusion.model.request.pod.PodStopRequest;
import com.luminous.fusion.model.response.agent.ContainerDto;
import com.luminous.fusion.model.response.pod.PodCreateResponse;
import com.luminous.fusion.service.PodService;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;
import java.util.Optional;

@RestController
@RequestMapping("/pods")
@AllArgsConstructor
@Slf4j
public class PodsController {

    private final PodService podService;

    @PostMapping("/")
    public ResponseEntity<PodCreateResponse> createPod(@RequestBody PodCreateRequest podCreateRequest) throws InterruptedException {
        log.info("Controller | createPod | PodCreateRequest : {}", podCreateRequest);

        String containerId = this.podService.createPod(podCreateRequest);

        return ResponseEntity.ok(
                new PodCreateResponse(containerId)
        );
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
    public ResponseEntity<ContainerDto> stopPod(@RequestBody PodStopRequest podStopRequest) throws InterruptedException {
        this.podService.stopPod(podStopRequest);
        Thread.sleep(1000);
        return ResponseEntity.ok(
                new ContainerDto(
                        this.podService.getSimpleContainerViaId(podStopRequest.getContainerId())
                )
        );
    }

    @PostMapping("/start")
    public ResponseEntity<ContainerDto> startPod(@RequestBody PodStartRequest podStartRequest) throws InterruptedException {
        this.podService.startPod(podStartRequest);
        Thread.sleep(1000);
        return ResponseEntity.ok(
                new ContainerDto(
                        this.podService.getSimpleContainerViaId(podStartRequest.getContainerId())
                )
        );
    }

    @GetMapping("/logs/{containerId}")
    public ResponseEntity<String> getPodLogs(@PathVariable String containerId) throws InterruptedException {
        return ResponseEntity.ok(
                this.podService.getContainerLogs(containerId)
        );
    }

}
