# Fusion

## Fusion Variables
This is variables that can be access via {{FUSION_VARIABLE_NAME}} in strings in the UI or other places. 
The token will be replaced with the actual configured value on the server.

### Global Fusion Variables
These variables are available throughout the application. 
These variables are either coming from the config file, or from generated values.

**Available Scopes**
- {{fusion.directory.root}} = Will give the current configured root directory for fusion

### Scoped Fusion Variables
These variables are avaible only on specific scopes. 

**Available Scopes**

- _Pod Creation Scopes_
  {{fusion.pod.id}}


### Custon Fusion Variables
In the config, custom fusion variables that can be set that can be access via {{fusion.custom._VALUE_}}. 