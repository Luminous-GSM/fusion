package com.luminous.fusion.model.domain.database;

import com.luminous.fusion.model.domain.types.ServerManagementStatusType;
import com.luminous.fusion.model.domain.types.UserPermissionType;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.List;
import java.util.Set;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class DatabaseModel {

    private List<UserDatabaseModal> users;

    private ServerManagementStatusType managementStatus;

    @AllArgsConstructor
    @NoArgsConstructor
    @Getter
    @Setter
    public static class UserDatabaseModal {
        private String username;
        private String email;
        private String password;
        private Set<UserPermissionType> permissions;
    }
}
