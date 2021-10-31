package com.luminous.fusion.security;

import com.luminous.fusion.configuration.LuminousPropertiesConfiguration;
import com.luminous.fusion.model.exception.InvalidAccessTokenException;
import lombok.AllArgsConstructor;
import org.springframework.security.core.authority.AuthorityUtils;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.stereotype.Component;

import javax.servlet.*;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;

@AllArgsConstructor
public class ApiKeyAuthenticationFilter implements Filter {

    private final LuminousPropertiesConfiguration luminousPropertiesConfiguration;

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain) throws IOException, ServletException {

        if(request instanceof HttpServletRequest && response instanceof HttpServletResponse) {

            String apiKey = ((HttpServletRequest) request).getHeader("x-api-key");
            if(apiKey != null) {

                if(apiKey.equals(luminousPropertiesConfiguration.getNode().getAuthorizationToken())) {
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
