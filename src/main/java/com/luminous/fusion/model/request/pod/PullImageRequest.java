package com.luminous.fusion.model.request.pod;

import lombok.*;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
@ToString
public class PullImageRequest {
    private String imageName;
    private String imageTag;
}
