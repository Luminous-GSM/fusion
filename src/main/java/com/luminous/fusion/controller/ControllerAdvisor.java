package com.luminous.fusion.controller;

import com.github.dockerjava.api.exception.ConflictException;
import com.github.dockerjava.api.exception.NotFoundException;
import com.github.dockerjava.api.exception.NotModifiedException;
import com.luminous.fusion.model.exception.ApiError;
import com.luminous.fusion.model.exception.UserNotFoundException;
import org.pf4j.PluginRuntimeException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseStatus;

import java.io.IOException;
import java.time.LocalDateTime;

@ControllerAdvice
public class ControllerAdvisor {

    @ExceptionHandler(PluginRuntimeException.class)
    public ResponseEntity<ApiError> handlePluginExceptions(PluginRuntimeException e) {
        return ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(
                        ApiError.builder()
                                .timestamp(LocalDateTime.now())
                                .exception(e.getClass().getName())
                                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                                .message("Plugin Manager error")
                                .debugMessage(e.getLocalizedMessage())
                                .build()
                );
    }

    @ExceptionHandler(IOException.class)
    public ResponseEntity<ApiError> handleIOExceptions(PluginRuntimeException e) {
        return ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(
                        ApiError.builder()
                                .timestamp(LocalDateTime.now())
                                .exception(e.getClass().getName())
                                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                                .message("IO Exception")
                                .debugMessage(e.getLocalizedMessage())
                                .build()
                );
    }

    @ExceptionHandler({UserNotFoundException.class})
    public ResponseEntity<ApiError> handleAuthenticationException(Exception e) {
        return ResponseEntity
                .status(HttpStatus.UNAUTHORIZED)
                .body(
                        ApiError.builder()
                                .timestamp(LocalDateTime.now())
                                .exception(e.getClass().getName())
                                .status(HttpStatus.UNAUTHORIZED)
                                .message("Not authorized")
                                .debugMessage(e.getLocalizedMessage())
                                .build()
                );
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<ApiError> handleException(Exception e) {
        return ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(
                        ApiError.builder()
                                .timestamp(LocalDateTime.now())
                                .exception(e.getClass().getName())
                                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                                .message("Something went wrong")
                                .debugMessage(e.getLocalizedMessage())
                                .build()
                );
    }

    @ExceptionHandler({NotFoundException.class, NotModifiedException.class, ConflictException.class})
    public ResponseEntity<ApiError> HandlePodException(Exception e) {
        return ResponseEntity
                .status(HttpStatus.BAD_REQUEST)
                .body(
                        ApiError.builder()
                                .timestamp(LocalDateTime.now())
                                .exception(e.getClass().getName())
                                .status(HttpStatus.BAD_REQUEST)
                                .message("Docker Pod error")
                                .debugMessage(e.getLocalizedMessage())
                                .build()
                );
    }



}
