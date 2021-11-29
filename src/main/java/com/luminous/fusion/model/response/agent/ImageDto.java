package com.luminous.fusion.model.response.agent;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.github.dockerjava.api.model.Image;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.Map;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class ImageDto {
    private Long created;
    private String id;
    private String parentId;
    private String[] repoTags;
    private String[] repoDigests;
    private Long size;
    private Long virtualSize;
    private Long sharedSize;
    public Map<String, String> labels;
    private Integer containers;

    public ImageDto(Image image) {
        this.created = image.getCreated();
        this.id = image.getId();
        this.parentId = image.getParentId();
        this.repoTags = image.getRepoTags();
        this.repoDigests = image.getRepoDigests();
        this.size = image.getSize();
        this.virtualSize = image.getVirtualSize();
        this.sharedSize = image.getSharedSize();
        this.labels = image.getLabels();
        this.containers = image.getContainers();
    }

}
