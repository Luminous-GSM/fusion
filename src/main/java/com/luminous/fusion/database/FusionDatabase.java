package com.luminous.fusion.database;

import com.luminous.fusion.model.domain.database.DatabaseModel;

import java.io.IOException;

public interface FusionDatabase {
    DatabaseModel readFile() throws IOException;
    void saveToFile(DatabaseModel databaseModel) throws IOException;
}
