package com.luminous.fusion.model.response.management;

import com.luminous.fusion.model.domain.types.ServerManagementStatusType;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class ManagementPingResult {
    private int userCount;
    private ServerManagementStatusType managementStatus;
}
