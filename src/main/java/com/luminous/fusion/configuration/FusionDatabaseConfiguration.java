package com.luminous.fusion.configuration;

import com.luminous.fusion.database.FusionDatabase;
import com.luminous.fusion.database.FusionDatabaseImpl;
import lombok.AllArgsConstructor;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class FusionDatabaseConfiguration {

    @Bean
    public FusionDatabase initializeFusionDatabase() {
        return new FusionDatabaseImpl();
    }

}
