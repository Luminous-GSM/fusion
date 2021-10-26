package com.luminous.fusion.service;

import com.luminous.fusion.database.FusionDatabase;
import com.luminous.fusion.model.domain.database.DatabaseModel;
import com.luminous.fusion.model.exception.UserNotFoundException;
import com.luminous.fusion.model.request.management.SignInRequest;
import com.luminous.fusion.model.response.management.ManagementPingResult;
import com.luminous.fusion.model.response.management.SignInResponse;
import lombok.AllArgsConstructor;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.util.Optional;

@Service
@AllArgsConstructor
public class DatabaseService {

    private final FusionDatabase database;

    public ManagementPingResult getStatus() throws IOException {
        DatabaseModel databaseModel = this.database.readFile();
        return new ManagementPingResult(databaseModel.getUsers().size(), databaseModel.getManagementStatus());
    }

    public SignInResponse authenticateUser(SignInRequest signInRequest) throws IOException {
        DatabaseModel databaseModel = this.database.readFile();
        Optional<DatabaseModel.UserDatabaseModal> optionalUserDatabaseModal = databaseModel
                .getUsers()
                .stream()
                .filter(userDatabaseModal ->
                        userDatabaseModal.getUsername().equals(signInRequest.getUsername())
                                && userDatabaseModal.getPassword().equals(signInRequest.getPassword())
                )
                .findFirst();

        if (optionalUserDatabaseModal.isPresent()) {
            return new SignInResponse("TOKEN");
        } else {
            throw new UserNotFoundException("User not found");
        }

    }

}
