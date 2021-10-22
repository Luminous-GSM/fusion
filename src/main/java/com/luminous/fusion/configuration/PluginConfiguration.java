package com.luminous.fusion.configuration;

import org.pf4j.DefaultExtensionFinder;
import org.pf4j.DefaultPluginManager;
import org.pf4j.ExtensionFinder;
import org.pf4j.PluginManager;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class PluginConfiguration {

    @Bean
    public PluginManager pluginManager() {
        PluginManager pluginManager = new DefaultPluginManager() {
            protected ExtensionFinder createExtensionFinder() {
                DefaultExtensionFinder extensionFinder = (DefaultExtensionFinder) super.createExtensionFinder();
                extensionFinder.addServiceProviderExtensionFinder();

                return extensionFinder;
            }
        };
        pluginManager.loadPlugins();
        pluginManager.startPlugins();

        return pluginManager;
    }

}
