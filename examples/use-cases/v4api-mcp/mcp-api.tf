resource "apim_apiv4" "mcp" {
  # should match the resource name
  hrid            = "message"
  name            = "[Terraform] Open Weather API With MCP"
  description     = "MCP API uses a API Open-Weather as tool."
  version         = "1,0"
  type            = "PROXY"
  state           = "STARTED"
  visibility      = "PRIVATE"
  lifecycle_state = "UNPUBLISHED"
  listeners = [
    {
      http = {
        type = "HTTP"
        entrypoints = [
          {
            type = "http-proxy"
          },
          {
            type = "mcp",
            configuration = jsonencode({
              mcpPath = "/mcp"
              tools = [
                {
                  gatewayMapping = {
                    http = {
                      method = "GET"
                      path   = "/forecast"
                      queryParams = [
                        "latitude",
                        "longitude",
                        "current",
                        "timezone",
                        "wind_speed_unit",
                      ]
                    }
                  }
                  toolDefinition = {
                    description = "Get current weather"
                    inputSchema = {
                      properties = {
                        current = {
                          description = "Comma-separated list of current weather variables to fetch."
                          example     = "temperature_2m,wind_speed_10m"
                          type        = "string"
                        }
                        latitude = {
                          description = "Latitude of the location."
                          example     = 40.7128
                          format      = "float"
                          type        = "number"
                        }
                        longitude = {
                          description = "Longitude of the location."
                          example     = -74.006
                          format      = "float"
                          type        = "number"
                        }
                        timezone = {
                          description = "Timezone for the output. Use `auto` to detect automatically."
                          example     = "auto"
                          type        = "string"
                        }
                        wind_speed_unit = {
                          description = "Unit for wind speed."
                          enum = [
                            "kmh",
                            "ms",
                            "mph",
                            "kn",
                          ]
                          example = "kmh"
                          type    = "string"
                        }
                      }
                      required = [
                        "latitude",
                        "longitude",
                        "current",
                        "timezone",
                        "wind_speed_unit",
                      ]
                      type = "object"
                    }
                    name = "get_forecast"
                  }
                },
              ]
            })
          }
        ]
        paths = [
          {
            path = "/open-weatherapis/"
          }
        ]
      }
    }
  ]
  endpoint_groups = [
    {
      name = "Default endpoint group"
      type = "http-proxy"
      endpoints = [
        {
          name   = "default"
          type   = "http-proxy"
          weight = 1
          configuration = jsonencode({
            target = "https://api.open-meteo.com/v1"
          })
        }
      ]
    }
  ]
  flow_execution = {
    match_required = false
    mode           = "DEFAULT"
  }
  flows = [
    {
      enabled = true
      selectors = [
        {
          type = "HTTP"
          "http" : {
            type         = "HTTP"
            path         = "/v1/forecast"
            pathOperator = "EQUALS"
            methods      = ["GET"]
          }
        }
      ]
      request = [
        {
          name = "mock"
          description : "7 day weather forecast for coordinates"
          enabled = true
          policy  = "mock"
          configuration = jsonencode({
            headers = [
              {
                name  = "Content-Type"
                value = "application/json"
              }
            ]
            status  = "200"
            content = <<-EOT
{
  "elevation": "44.812",
  "hourly_units": {
    "additionalProperty": "Mocked string"
  },
  "generationtime_ms": "2.2119",
  "daily_units": {
    "additionalProperty": "Mocked string"
  },
  "latitude": "52.52",
  "daily": {
    "shortwave_radiation_sum": [
      "0.9311039678493772"
    ],
    "sunrise": [
      "0.8960258575716086"
    ],
    "wind_speed_10m_max": [
      "0.23497278680038025"
    ],
    "apparent_temperature_min": [
      "0.962365342792647"
    ],
    "uv_index_max": [
      "0.21604074164840126"
    ],
    "temperature_2m_min": [
      "0.9146276316046636"
    ],
    "et0_fao_evapotranspiration": [
      "0.045100544361300954"
    ],
    "wind_gusts_10m_max": [
      "0.9187254684638405"
    ],
    "uv_index_clear_sky_max": [
      "0.884193403221878"
    ],
    "apparent_temperature_max": [
      "0.5204369223251467"
    ],
    "sunset": [
      "0.11586878455679994"
    ],
    "temperature_2m_max": [
      "0.8871327060740304"
    ],
    "wind_direction_10m_dominant": [
      "0.8907798772493127"
    ],
    "time": [
      "Mocked string"
    ],
    "weather_code": [
      "0.4258971242744911"
    ],
    "precipitation_sum": [
      "0.48366701987049443"
    ],
    "precipitation_hours": [
      "0.4610770379807593"
    ]
  },
  "utc_offset_seconds": "3600",
  "hourly": {
    "soil_moisture_27_81cm": [
      "0.2767567625200661"
    ],
    "relative_humidity_2m": [
      "0.7116329549657038"
    ],
    "cloud_cover_low": [
      "0.46082254602906325"
    ],
    "vapour_pressure_deficit": [
      "0.8191199623791499"
    ],
    "soil_temperature_0cm": [
      "0.5117555108353147"
    ],
    "soil_moisture_9_27cm": [
      "0.9186972072799153"
    ],
    "evapotranspiration": [
      "0.32939357463854857"
    ],
    "soil_temperature_6cm": [
      "0.4625992240963396"
    ],
    "shortwave_radiation": [
      "0.7136354250090723"
    ],
    "precipitation": [
      "0.20188187855605766"
    ],
    "soil_temperature_54cm": [
      "0.5624147350667345"
    ],
    "direct_radiation": [
      "0.6591336200309232"
    ],
    "soil_moisture_0_1cm": [
      "0.8202633573791254"
    ],
    "cloud_cover": [
      "0.049191248967149326"
    ],
    "apparent_temperature": [
      "0.4602027698500172"
    ],
    "wind_speed_80m": [
      "0.8249629078386708"
    ],
    "cloud_cover_high": [
      "0.5661366602124245"
    ],
    "soil_moisture_3_9cm": [
      "0.4085153046306096"
    ],
    "soil_temperature_18cm": [
      "0.1632756238129799"
    ],
    "pressure_msl": [
      "0.3252521144774868"
    ],
    "wind_speed_10m": [
      "0.4084018888187052"
    ],
    "dew_point_2m": [
      "0.83659642462144"
    ],
    "wind_speed_180m": [
      "0.26677960002851764"
    ],
    "wind_direction_80m": [
      "0.5235938083967235"
    ],
    "snow_height": [
      "0.8752803218729202"
    ],
    "wind_direction_10m": [
      "0.2843038988373786"
    ],
    "wind_gusts_10m": [
      "0.1683774784578721"
    ],
    "freezing_level_height": [
      "0.5467270998515629"
    ],
    "temperature_2m": [
      "0.18804468405669283"
    ],
    "cloud_cover_mid": [
      "0.94837371946797"
    ],
    "wind_direction_120m": [
      "0.7703364175102"
    ],
    "wind_direction_180m": [
      "0.11455530728856145"
    ],
    "soil_moisture_1_3cm": [
      "0.5405602834735213"
    ],
    "time": [
      "Mocked string"
    ],
    "wind_speed_120m": [
      "0.9351211913437067"
    ],
    "direct_normal_irradiance": [
      "0.5701448431404776"
    ],
    "diffuse_radiation": [
      "0.9879760182489068"
    ],
    "weather_code": [
      "0.2054339799667142"
    ]
  },
  "current_weather": {
    "temperature": "0.387407897329341",
    "wind_speed": "0.4370495943911785",
    "wind_direction": "0.9225714265890317",
    "time": "Mocked string",
    "weather_code": "343"
  },
  "longitude": "0.20848379195333233"
}
EOT
          })
        }
      ]
    }
  ]
  analytics = {
    enabled = true
  }
  plans = [
    {
      hrid        = "key-less"
      name        = "No security"
      type        = "API"
      mode        = "STANDARD"
      validation  = "AUTO"
      status      = "PUBLISHED"
      description = "This plan does not require any authentication"
      security = {
        type = "KEY_LESS"
      }
    }
  ]
  pages = [
    {
      hrid       = "aside"
      name       = "Aside"
      type       = "SYSTEM_FOLDER"
      order      = 0
      published  = true
      visibility = "PUBLIC"
      homepage   = false
    },
    {
      hrid       = "swagger"
      name       = "Swagger"
      type       = "SWAGGER"
      order      = 1
      published  = true
      visibility = "PUBLIC"
      configuration = {
        viewer                 = "Swagger"
        entrypointAsBasePath   = "false"
        entrypointsAsServers   = "false"
        tryIt                  = "true"
        disableSyntaxHighlight = "false"
        tryItAnonymous         = "false"
        showURL                = "false"
        displayOperationId     = "false"
        usePkce                = "false"
        docExpansion           = "none"
        enableFiltering        = "false"
        showExtensions         = "false"
        showCommonExtensions   = "false"
        maxDisplayedTags       = "-1"
      }
      homepage = false
      content  = <<-EOT
openapi: 3.0.0
info:
  title: Open-Weather APIs
  description: 'Open-Meteo offers free weather forecast APIs for open-source developers and non-commercial use. No API key is required.'
  version: '1.0'
  contact:
    name: Open-Meteo
    url: https://open-meteo.com
    email: info@open-meteo.com
  license:
    name: Attribution 4.0 International (CC BY 4.0)
    url: https://creativecommons.org/licenses/by/4.0/
  termsOfService: https://open-meteo.com/en/features#terms
paths:
  /v1/forecast:
    servers:
      - url: https://api.open-meteo.com
    get:
      tags:
      - Weather Forecast APIs
      summary: 7 day weather forecast for coordinates
      description: 7 day weather variables in hourly and daily resolution for given WGS84 latitude and longitude coordinates. Available worldwide.
      parameters:
      - name: hourly
        in: query
        explode: false
        schema:
          type: array
          items:
            type: string
            enum:
            - temperature_2m
            - relative_humidity_2m
            - dew_point_2m
            - apparent_temperature
            - pressure_msl
            - cloud_cover
            - cloud_cover_low
            - cloud_cover_mid
            - cloud_cover_high
            - wind_speed_10m
            - wind_speed_80m
            - wind_speed_120m
            - wind_speed_180m
            - wind_direction_10m
            - wind_direction_80m
            - wind_direction_120m
            - wind_direction_180m
            - wind_gusts_10m
            - shortwave_radiation
            - direct_radiation
            - direct_normal_irradiance
            - diffuse_radiation
            - vapour_pressure_deficit
            - evapotranspiration
            - precipitation
            - weather_code
            - snow_height
            - freezing_level_height
            - soil_temperature_0cm
            - soil_temperature_6cm
            - soil_temperature_18cm
            - soil_temperature_54cm
            - soil_moisture_0_1cm
            - soil_moisture_1_3cm
            - soil_moisture_3_9cm
            - soil_moisture_9_27cm
            - soil_moisture_27_81cm
      - name: daily
        in: query
        schema:
          type: array
          items:
            type: string
            enum:
            - temperature_2m_max
            - temperature_2m_min
            - apparent_temperature_max
            - apparent_temperature_min
            - precipitation_sum
            - precipitation_hours
            - weather_code
            - sunrise
            - sunset
            - wind_speed_10m_max
            - wind_gusts_10m_max
            - wind_direction_10m_dominant
            - shortwave_radiation_sum
            - uv_index_max
            - uv_index_clear_sky_max
            - et0_fao_evapotranspiration
      - name: latitude
        in: query
        required: true
        description: "WGS84 coordinate"
        schema:
          type: number
          format: double
      - name: longitude
        in: query
        required: true
        description: "WGS84 coordinate"
        schema:
          type: number
          format: double
      - name: current_weather
        in: query
        schema:
          type: boolean
      - name: temperature_unit
        in: query
        schema:
          type: string
          default: celsius
          enum:
          - celsius
          - fahrenheit
      - name: wind_speed_unit
        in: query
        schema:
          type: string
          default: kmh
          enum:
          - kmh
          - ms
          - mph
          - kn
      - name: timeformat
        in: query
        description: If format `unixtime` is selected, all time values are returned in UNIX epoch time in seconds. Please not that all time is then in GMT+0! For daily values with unix timestamp, please apply `utc_offset_seconds` again to get the correct date.
        schema:
          type: string
          default: iso8601
          enum:
          - iso8601
          - unixtime
      - name: timezone
        in: query
        description: If `timezone` is set, all timestamps are returned as local-time and data is returned starting at 0:00 local-time. Any time zone name from the [time zone database](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) is supported.
        schema:
          type: string
      - name: past_days
        in: query
        description: If `past_days` is set, yesterdays or the day before yesterdays data are also returned.
        schema:
          type: integer
          enum:
          - 1
          - 2
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  latitude:
                    type: number
                    example: 52.52
                    description: WGS84 of the center of the weather grid-cell which was used to generate this forecast. This coordinate might be up to 5 km away.
                  longitude:
                    type: number
                    example: 13.419.52
                    description: WGS84 of the center of the weather grid-cell which was used to generate this forecast. This coordinate might be up to 5 km away.
                  elevation:
                    type: number
                    example: 44.812
                    description: The elevation in meters of the selected weather grid-cell. In mountain terrain it might differ from the location you would expect.
                  generationtime_ms:
                    type: number
                    example: 2.2119
                    description: Generation time of the weather forecast in milli seconds. This is mainly used for performance monitoring and improvements.
                  utc_offset_seconds:
                    type: integer
                    example: 3600
                    description: Applied timezone offset from the &timezone= parameter.
                  hourly:
                    $ref: "#/components/schemas/HourlyResponse"
                  hourly_units:
                    type: object
                    additionalProperties:
                      type: string
                    description: For each selected weather variable, the unit will be listed here.
                  daily:
                    $ref: "#/components/schemas/DailyResponse"
                  daily_units:
                    type: object
                    additionalProperties:
                      type: string
                    description: For each selected daily weather variable, the unit will be listed here.
                  current_weather:
                    $ref: "#/components/schemas/CurrentWeather"
        "400":
          description: Bad Request
          content:
            application/json:
              schema:
                type: object
                properties:
                  error:
                    type: boolean
                    description: Always set true for errors
                  reason:
                    type: string
                    description: Description of the error
                    example: "Latitude must be in range of -90 to 90°. Given: 300"
components:
  schemas:
    HourlyResponse:
      type: object
      description: For each selected weather variable, data will be returned as a floating point array. Additionally a `time` array will be returned with ISO8601 timestamps.
      required:
        - time
      properties:
        time:
          type: array
          items:
            type: string
        temperature_2m:
          type: array
          items:
            type: number
        relative_humidity_2m:
          type: array
          items:
            type: number
        dew_point_2m:
          type: array
          items:
            type: number
        apparent_temperature:
          type: array
          items:
            type: number
        pressure_msl:
          type: array
          items:
            type: number
        cloud_cover:
          type: array
          items:
            type: number
        cloud_cover_low:
          type: array
          items:
            type: number
        cloud_cover_mid:
          type: array
          items:
            type: number
        cloud_cover_high:
          type: array
          items:
            type: number
        wind_speed_10m:
          type: array
          items:
            type: number
        wind_speed_80m:
          type: array
          items:
            type: number
        wind_speed_120m:
          type: array
          items:
            type: number
        wind_speed_180m:
          type: array
          items:
            type: number
        wind_direction_10m:
          type: array
          items:
            type: number
        wind_direction_80m:
          type: array
          items:
            type: number
        wind_direction_120m:
          type: array
          items:
            type: number
        wind_direction_180m:
          type: array
          items:
            type: number
        wind_gusts_10m:
          type: array
          items:
            type: number
        shortwave_radiation:
          type: array
          items:
            type: number
        direct_radiation:
          type: array
          items:
            type: number
        direct_normal_irradiance:
          type: array
          items:
            type: number
        diffuse_radiation:
          type: array
          items:
            type: number
        vapour_pressure_deficit:
          type: array
          items:
            type: number
        evapotranspiration:
          type: array
          items:
            type: number
        precipitation:
          type: array
          items:
            type: number
        weather_code:
          type: array
          items:
            type: number
        snow_height:
          type: array
          items:
            type: number
        freezing_level_height:
          type: array
          items:
            type: number
        soil_temperature_0cm:
          type: array
          items:
            type: number
        soil_temperature_6cm:
          type: array
          items:
            type: number
        soil_temperature_18cm:
          type: array
          items:
            type: number
        soil_temperature_54cm:
          type: array
          items:
            type: number
        soil_moisture_0_1cm:
          type: array
          items:
            type: number
        soil_moisture_1_3cm:
          type: array
          items:
            type: number
        soil_moisture_3_9cm:
          type: array
          items:
            type: number
        soil_moisture_9_27cm:
          type: array
          items:
            type: number
        soil_moisture_27_81cm:
          type: array
          items:
            type: number
    DailyResponse:
      type: object
      description: For each selected daily weather variable, data will be returned as a floating point array. Additionally a `time` array will be returned with ISO8601 timestamps.
      properties:
        time:
          type: array
          items:
            type: string
        temperature_2m_max:
          type: array
          items:
            type: number
        temperature_2m_min:
          type: array
          items:
            type: number
        apparent_temperature_max:
          type: array
          items:
            type: number
        apparent_temperature_min:
          type: array
          items:
            type: number
        precipitation_sum:
          type: array
          items:
            type: number
        precipitation_hours:
          type: array
          items:
            type: number
        weather_code:
          type: array
          items:
            type: number
        sunrise:
          type: array
          items:
            type: number
        sunset:
          type: array
          items:
            type: number
        wind_speed_10m_max:
          type: array
          items:
            type: number
        wind_gusts_10m_max:
          type: array
          items:
            type: number
        wind_direction_10m_dominant:
          type: array
          items:
            type: number
        shortwave_radiation_sum:
          type: array
          items:
            type: number
        uv_index_max:
          type: array
          items:
            type: number
        uv_index_clear_sky_max:
          type: array
          items:
            type: number
        et0_fao_evapotranspiration:
          type: array
          items:
            type: number
      required:
        - time
    CurrentWeather:
      type: object
      description: "Current weather conditions with the attributes: time, temperature, wind_speed, wind_direction and weather_code"
      properties:
        time:
          type: string
        temperature:
          type: number
        wind_speed:
          type: number
        wind_direction:
          type: number
        weather_code:
          type: integer
      required:
        - time
        - temperature
        - wind_speed
        - wind_direction
        - weather_code
EOT
  }]

}