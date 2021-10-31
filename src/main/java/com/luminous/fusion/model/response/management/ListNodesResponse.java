package com.luminous.fusion.model.response.management;

import com.luminous.fusion.model.domain.server.NodeDescription;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.List;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class ListNodesResponse {
    List<NodeDescription> nodes;
}
