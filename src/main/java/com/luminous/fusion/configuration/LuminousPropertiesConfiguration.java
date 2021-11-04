package com.luminous.fusion.configuration;

import com.luminous.fusion.model.domain.server.HostingPlatform;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

@Configuration
@ConfigurationProperties(prefix = "luminous")
@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class LuminousPropertiesConfiguration {

    private Node node = new Node("FusionServer", "Fusion server for hosting games", "AuthToken", "UniqueId");

    private String version;

    private HostingPlatform platform;

    private Docker docker;

    @AllArgsConstructor
    @NoArgsConstructor
    @Getter
    @Setter
    public static class Docker {
        private String host;
    }

    @AllArgsConstructor
    @NoArgsConstructor
    @Getter
    @Setter
    public static class Node {
        private String name;
        private String description;
        private String authorizationToken = "NONE";
        private String uniqueId = "NONE";
    }

}
