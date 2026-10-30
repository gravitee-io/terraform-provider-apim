resource "apim_portal_theme" "brand" {
  hrid = "brand-theme"
  name = "[Terraform] Brand theme"
  definition = {
    custom_css = "body { line-height: 1.5; }"
    font       = { font_family = "Arial, sans-serif" }
    color = {
      primary    = "#146b50"
      secondary  = "#344256"
      tertiary   = "#7a8899"
      error      = "#b3261e"
      background = { page = "#ffffff", card = "#f2f4f6" }
    }
  }
}
