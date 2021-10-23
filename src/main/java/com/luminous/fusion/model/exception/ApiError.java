package com.luminous.fusion.model.exception;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Getter;
import lombok.Setter;
import org.springframework.http.HttpStatus;

import java.time.LocalDateTime;
import java.util.List;

@Getter
@Setter
@Builder
@AllArgsConstructor
public class ApiError {
    private HttpStatus status;
    private LocalDateTime timestamp;
    private String message;
    private String debugMessage;
    private String exception;
    private List<ApiSubError> subErrors;

    public ApiError() {
        timestamp = LocalDateTime.now();
    }

}
