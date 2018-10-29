// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "title": "Robot Service",
    "version": "0.0.1"
  },
  "paths": {
    "/recommendation_plugins": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "Returns list of the installed recommendation plugins",
        "operationId": "getRecomendationPlugins",
        "responses": {
          "200": {
            "description": "no error",
            "schema": {
              "type": "object",
              "properties": {
                "InstalledRecommendationPlugins": {
                  "description": "installed plugins",
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/RecommendationPlugin"
                  }
                },
                "TotalCount": {
                  "type": "integer"
                }
              }
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "CheckResult": {
      "description": "CheckResult represents the single result of Check function invocation of specific plugin.",
      "type": "object",
      "properties": {
        "CheckStatus": {
          "description": "shows check status",
          "type": "string",
          "enum": [
            "RED",
            "YELLOW",
            "GREEN"
          ]
        },
        "CompletedAt": {
          "description": "date/Time of check execution",
          "type": "string",
          "format": "dateTime"
        },
        "Description": {
          "description": "detailed check result description",
          "type": "string"
        },
        "ExecutionStatus": {
          "description": "shows check execution errors",
          "type": "string"
        },
        "Id": {
          "description": "unique UUID of Check function invocation of specific plugin",
          "type": "string"
        },
        "Name": {
          "description": "check name",
          "type": "string"
        },
        "PossibleActions": {
          "description": "list of possible actions to fix caveats check was found",
          "type": "array",
          "items": {
            "$ref": "#/definitions/PluginAction"
          }
        }
      }
    },
    "PluginAction": {
      "description": "CheckResult represents the single result of Check function invocation of specific plugin.",
      "type": "object",
      "properties": {
        "Description": {
          "description": "detailed action description",
          "type": "string"
        },
        "Id": {
          "description": "unique UUID of plugin action",
          "type": "string"
        }
      }
    },
    "RecommendationPlugin": {
      "description": "RecommendationPlugin represents the installed recommendation plugin",
      "type": "object",
      "properties": {
        "Description": {
          "description": "detailed plugin description",
          "type": "string"
        },
        "Id": {
          "description": "unique ID of installed plugin\nbasically it is slugged URI of plugin repository name e. g. supergiant-request-limits-check\n",
          "type": "string"
        },
        "InstalledAt": {
          "description": "date/Time the plugin was installed",
          "type": "string",
          "format": "dateTime"
        },
        "Name": {
          "description": "name is the name of the plugin.",
          "type": "string"
        },
        "Status": {
          "description": "plugin status",
          "type": "string"
        },
        "Version": {
          "description": "plugin version, major version shall be equal to robots version",
          "type": "string"
        }
      }
    },
    "error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "format": "int64"
        },
        "message": {
          "type": "string"
        }
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "title": "Robot Service",
    "version": "0.0.1"
  },
  "paths": {
    "/recommendation_plugins": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "Returns list of the installed recommendation plugins",
        "operationId": "getRecomendationPlugins",
        "responses": {
          "200": {
            "description": "no error",
            "schema": {
              "type": "object",
              "properties": {
                "InstalledRecommendationPlugins": {
                  "description": "installed plugins",
                  "type": "array",
                  "items": {
                    "$ref": "#/definitions/RecommendationPlugin"
                  }
                },
                "TotalCount": {
                  "type": "integer"
                }
              }
            }
          },
          "default": {
            "description": "error",
            "schema": {
              "$ref": "#/definitions/error"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "CheckResult": {
      "description": "CheckResult represents the single result of Check function invocation of specific plugin.",
      "type": "object",
      "properties": {
        "CheckStatus": {
          "description": "shows check status",
          "type": "string",
          "enum": [
            "RED",
            "YELLOW",
            "GREEN"
          ]
        },
        "CompletedAt": {
          "description": "date/Time of check execution",
          "type": "string",
          "format": "dateTime"
        },
        "Description": {
          "description": "detailed check result description",
          "type": "string"
        },
        "ExecutionStatus": {
          "description": "shows check execution errors",
          "type": "string"
        },
        "Id": {
          "description": "unique UUID of Check function invocation of specific plugin",
          "type": "string"
        },
        "Name": {
          "description": "check name",
          "type": "string"
        },
        "PossibleActions": {
          "description": "list of possible actions to fix caveats check was found",
          "type": "array",
          "items": {
            "$ref": "#/definitions/PluginAction"
          }
        }
      }
    },
    "PluginAction": {
      "description": "CheckResult represents the single result of Check function invocation of specific plugin.",
      "type": "object",
      "properties": {
        "Description": {
          "description": "detailed action description",
          "type": "string"
        },
        "Id": {
          "description": "unique UUID of plugin action",
          "type": "string"
        }
      }
    },
    "RecommendationPlugin": {
      "description": "RecommendationPlugin represents the installed recommendation plugin",
      "type": "object",
      "properties": {
        "Description": {
          "description": "detailed plugin description",
          "type": "string"
        },
        "Id": {
          "description": "unique ID of installed plugin\nbasically it is slugged URI of plugin repository name e. g. supergiant-request-limits-check\n",
          "type": "string"
        },
        "InstalledAt": {
          "description": "date/Time the plugin was installed",
          "type": "string",
          "format": "dateTime"
        },
        "Name": {
          "description": "name is the name of the plugin.",
          "type": "string"
        },
        "Status": {
          "description": "plugin status",
          "type": "string"
        },
        "Version": {
          "description": "plugin version, major version shall be equal to robots version",
          "type": "string"
        }
      }
    },
    "error": {
      "type": "object",
      "required": [
        "message"
      ],
      "properties": {
        "code": {
          "type": "integer",
          "format": "int64"
        },
        "message": {
          "type": "string"
        }
      }
    }
  }
}`))
}
