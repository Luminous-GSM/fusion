package com.luminous.fusion.model.exception;

import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.ResponseStatus;

@ResponseStatus(code = HttpStatus.UNAUTHORIZED, reason = "Invalid access token provided")
public class InvalidAccessTokenException extends RuntimeException {
    public InvalidAccessTokenException() {
    }

    public InvalidAccessTokenException(String message) {
        super(message);
    }

    public InvalidAccessTokenException(String message, Throwable cause) {
        super(message, cause);
    }

    public InvalidAccessTokenException(Throwable cause) {
        super(cause);
    }

    public InvalidAccessTokenException(String message, Throwable cause, boolean enableSuppression, boolean writableStackTrace) {
        super(message, cause, enableSuppression, writableStackTrace);
    }
}
