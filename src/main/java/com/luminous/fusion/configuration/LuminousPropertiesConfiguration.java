package com.luminous.fusion.configuration;

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

    private String version;

    // Default value makes default token optional via OS environment
    private String defaultToken = "NONE";

    private String platform;

    private Database database;

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
    public static class Database {
        private String filePath;
    }

}
