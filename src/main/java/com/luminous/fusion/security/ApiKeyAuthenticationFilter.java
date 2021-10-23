package com.luminous.fusion.security;

import com.luminous.fusion.model.exception.InvalidAccessTokenException;
import org.springframework.security.core.authority.AuthorityUtils;
import org.springframework.security.core.context.SecurityContextHolder;

import javax.servlet.*;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

public class ApiKeyAuthenticationFilter implements Filter {

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain) throws IOException, ServletException {

        if(request instanceof HttpServletRequest && response instanceof HttpServletResponse) {

            String apiKey = ((HttpServletRequest) request).getHeader("x-api-key");
            if(apiKey != null) {

                // TODO Check the api key here if it's valid. Should be a valid JWT.
                boolean valid = true;
                if(valid) {
                    // TODO Get authorities from JWT token
                    ApiKeyAuthenticationToken apiToken = new ApiKeyAuthenticationToken(apiKey, AuthorityUtils.NO_AUTHORITIES);
                    SecurityContextHolder.getContext().setAuthentication(apiToken);
                } else {
                    throw new InvalidAccessTokenException("Invalid API Key");
//                    HttpServletResponse httpResponse = (HttpServletResponse) response;
//                    httpResponse.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
//                    httpResponse.getWriter().write("Invalid API Key");
//                    return;
                }
            }
        }

        chain.doFilter(request, response);

    }
}
