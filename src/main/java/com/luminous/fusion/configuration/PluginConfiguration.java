package com.luminous.fusion.configuration;

import org.pf4j.*;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class PluginConfiguration {

    @Bean
    public PluginManager pluginManager() {
//        PluginManager pluginManager = new DefaultPluginManager() {
//            protected ExtensionFinder createExtensionFinder() {
//                DefaultExtensionFinder extensionFinder = (DefaultExtensionFinder) super.createExtensionFinder();
//                extensionFinder.addServiceProviderExtensionFinder();
//
//                return extensionFinder;
//            }
//        };
        PluginManager pluginManager = new JarPluginManager();
        pluginManager.loadPlugins();
        pluginManager.startPlugins();

        return pluginManager;
    }

}
