package com.luminous.fusion.database;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.luminous.fusion.configuration.LuminousPropertiesConfiguration;
import com.luminous.fusion.model.domain.database.DatabaseModel;
import com.luminous.fusion.model.domain.types.ServerManagementStatusType;
import com.luminous.fusion.model.domain.types.UserPermissionType;
import lombok.AllArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.io.Resource;
import org.springframework.stereotype.Component;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.util.List;
import java.util.Set;

@Slf4j
public class FusionDatabaseImpl implements FusionDatabase {

    @Value("classpath:data/database.json")
    Resource databaseFile;

    private final ObjectMapper objectMapper = new ObjectMapper();

    @Override
    public DatabaseModel readFile() throws IOException {
        return objectMapper.readValue(databaseFile.getInputStream(), DatabaseModel.class);
    }

    @Override
    public void saveToFile(DatabaseModel databaseModel) throws IOException {
        objectMapper.writeValue(databaseFile.getFile(), databaseModel);
    }
}
