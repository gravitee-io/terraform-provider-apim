resource "apim_portal_theme" "brand" {
  hrid = "brand-theme"
  name = "Brand theme"
  definition = {
    font = { font_family = "Arial, sans-serif" }
    color = {
      primary    = "#146b50"
      secondary  = "#344256"
      background = { page = "#ffffff", card = "#f2f4f6" }
    }
  }
}

resource "apim_portal" "developer-portal" {
  hrid              = "developer-portal"
  name              = "Developer Portal"
  structure         = {}
  active_theme_hrid = apim_portal_theme.brand.hrid
}
