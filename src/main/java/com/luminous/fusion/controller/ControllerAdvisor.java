package com.luminous.fusion.controller;

import com.github.dockerjava.api.exception.DockerException;
import com.luminous.fusion.model.exception.ApiError;
import com.luminous.fusion.model.exception.InvalidAccessTokenException;
import com.luminous.fusion.model.exception.UserNotFoundException;
import lombok.extern.slf4j.Slf4j;
import org.pf4j.PluginRuntimeException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;

import java.io.IOException;
import java.time.LocalDateTime;

@ControllerAdvice
@Slf4j
public class ControllerAdvisor {

    @ExceptionHandler(PluginRuntimeException.class)
    public ResponseEntity<ApiError> handlePluginExceptions(PluginRuntimeException e) {
        ApiError error = ApiError.builder()
                .timestamp(LocalDateTime.now())
                .exception(e.getClass().getName())
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .message("Plugin Manager error")
                .debugMessage(e.getLocalizedMessage())
                .build();

        log.error("Plugin Manager error", e);
        log.error("{}", error);

        return ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(error);
    }

    @ExceptionHandler(IOException.class)
    public ResponseEntity<ApiError> handleIOExceptions(PluginRuntimeException e) {
        ApiError error = ApiError.builder()
                .timestamp(LocalDateTime.now())
                .exception(e.getClass().getName())
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .message("IO Exception")
                .debugMessage(e.getLocalizedMessage())
                .build();

        log.error("IO Exception", e);
        log.error("{}", error);

        return ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(error);
    }

    @ExceptionHandler({UserNotFoundException.class, InvalidAccessTokenException.class})
    public ResponseEntity<ApiError> handleAuthenticationException(Exception e) {
        ApiError error = ApiError.builder()
                .timestamp(LocalDateTime.now())
                .exception(e.getClass().getName())
                .status(HttpStatus.UNAUTHORIZED)
                .message("Not authorized")
                .debugMessage(e.getLocalizedMessage())
                .build();

        log.error("Not authorized", e);
        log.error("{}", error);

        return ResponseEntity
                .status(HttpStatus.UNAUTHORIZED)
                .body(error);
    }



    @ExceptionHandler({DockerException.class})
    public ResponseEntity<ApiError> HandlePodException(DockerException e) {
        ApiError error = ApiError.builder()
                .timestamp(LocalDateTime.now())
                .exception(DockerException.class.getName())
                .status(HttpStatus.resolve(e.getHttpStatus()))
                .message("Fusion pod error")
                .debugMessage(
                        e.getMessage().replace(
                                String.format("Status %d: ", e.getHttpStatus()),
                                ""
                        )
                )
                .build();

        log.error("Fusion pod error", e);
        log.error("{}", error);

        return ResponseEntity
                .status(HttpStatus.BAD_REQUEST)
                .body(error);
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<ApiError> handleException(Exception e) {
        ApiError error = ApiError.builder()
                .timestamp(LocalDateTime.now())
                .exception(e.getClass().getName())
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .message("Something went wrong")
                .debugMessage(e.getMessage())
                .build();

        log.error("Something went wrong", e);
        log.error("{}", error);

        return ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(error);
    }



}
