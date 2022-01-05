package com.luminous.fusion.configuration;

import com.luminous.fusion.model.domain.server.HostingPlatform;
import lombok.*;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

@Configuration
@ConfigurationProperties(prefix = "luminous")
@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
@ToString
public class LuminousPropertiesConfiguration {

    private Node node = new Node("FusionServer", "Fusion server for hosting games", "AuthToken", "UniqueId", "fusion-node-XXXXX");

    private String version;

    private HostingPlatform platform;

    @AllArgsConstructor
    @NoArgsConstructor
    @Getter
    @Setter
    @ToString
    public static class Node {
        private String name;
        private String description;
        private String authorizationToken = "NONE";
        private String uniqueId = "NONE";
        private String hostname = "fusion-node-XXXXX";
    }
}
