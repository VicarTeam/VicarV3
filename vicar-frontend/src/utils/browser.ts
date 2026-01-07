type BrowserName = "Chrome"|"Firefox"|"Edge"|"Safari"|"Opera"|"IE"|"Chromium"|"Unknown Browser";

export function getBrowserName(): BrowserName {
  const ua = navigator.userAgent;
  if (ua.includes("Chrome")) {
    return "Chrome";
  } else if (ua.includes("Firefox")) {
    return "Firefox";
  } else if (ua.includes("Edg")) {
    return "Edge";
  } else if (ua.includes("Safari")) {
    return "Safari";
  } else if (ua.includes("OPR")) {
    return "Opera";
  } else if (ua.includes("Chromium")) {
    return "Chromium";
  } else if (ua.includes("MSIE") || ua.includes("Trident")) {
    return "IE";
  } else {
    return "Unknown Browser";
  }
}

export function getBrowserFaIcon(name: BrowserName): string {
  switch (name) {
    case "Chrome":
    case "Chromium":
      return "fab fa-chrome";
    case "Firefox":
      return "fab fa-firefox";
    case "Edge":
      return "fab fa-edge";
    case "Safari":
      return "fab fa-safari";
    case "Opera":
      return "fab fa-opera";
    case "IE":
      return "fab fa-internet-explorer";
    default:
      return "fas fa-question";
  }
}