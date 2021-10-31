package com.luminous.fusion.model.domain.server;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class NodeDescription {
    String nodeUniqueId;
    String name;
    String description;
    NodeStatus nodeStatus;
}
