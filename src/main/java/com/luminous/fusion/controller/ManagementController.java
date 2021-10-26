package com.luminous.fusion.controller;

import com.luminous.fusion.model.request.management.SignInRequest;
import com.luminous.fusion.model.response.management.ManagementPingResult;
import com.luminous.fusion.model.response.management.SignInResponse;
import com.luminous.fusion.service.DatabaseService;
import lombok.AllArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.io.IOException;

@RestController
@RequestMapping("/management")
@AllArgsConstructor
public class ManagementController {

    private final DatabaseService databaseService;

    @GetMapping("/ping")
    public ResponseEntity<ManagementPingResult> pingPong() throws IOException {
        return ResponseEntity.ok(this.databaseService.getStatus());
    }

    @PostMapping("/sign-in")
    public ResponseEntity<SignInResponse> authorizeUser(@RequestBody SignInRequest signInRequest) throws IOException {
        return ResponseEntity.ok(this.databaseService.authenticateUser(signInRequest));
    }

}
