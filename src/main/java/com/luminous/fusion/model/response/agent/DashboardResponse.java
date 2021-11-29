package com.luminous.fusion.model.response.agent;

import com.github.dockerjava.api.model.Container;
import com.github.dockerjava.api.model.Image;
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
public class DashboardResponse {
    private NodeDescription nodeDescription;
    private List<ContainerDto> pods;
    private List<ImageDto> images;

}
