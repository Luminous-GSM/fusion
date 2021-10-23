package com.luminous.fusion.model.request.pod;

import com.luminous.fusion.model.domain.pod.PodDescription;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
public class PodCreateRequest {
    private PodDescription podDescription;
}
