package com.luminous.fusion.model.request.pod;

import com.luminous.fusion.model.domain.pod.PodDescription;
import lombok.*;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
@ToString
public class PodCreateRequest {
    private PodDescription podDescription;
}
