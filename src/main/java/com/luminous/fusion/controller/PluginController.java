package com.luminous.fusion.controller;

import com.luminous.fusion.service.PluginService;
import lombok.AllArgsConstructor;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/plugin")
@AllArgsConstructor
public class PluginController {

    private final PluginService pluginService;



}
