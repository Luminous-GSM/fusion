package com.luminous.fusion;

import com.luminous.fusion.configuration.LuminousPropertiesConfiguration;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.context.properties.EnableConfigurationProperties;

@SpringBootApplication
@EnableConfigurationProperties(LuminousPropertiesConfiguration.class)
public class FusionApplication {

	public static void main(String[] args) {
		SpringApplication.run(FusionApplication.class, args);
	}

}
