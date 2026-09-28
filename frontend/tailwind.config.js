function colorVar(name) {
  return `rgb(var(--color-${name}) / <alpha-value>)`
}

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx,html}",
  ],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        "surface-container-high": colorVar("surface-container-high"),
        "primary": colorVar("primary"),
        "on-primary-fixed": colorVar("on-primary-fixed"),
        "surface-variant": colorVar("surface-variant"),
        "surface": colorVar("surface"),
        "on-tertiary-fixed-variant": colorVar("on-tertiary-fixed-variant"),
        "tertiary-fixed": colorVar("tertiary-fixed"),
        "inverse-on-surface": colorVar("inverse-on-surface"),
        "surface-bright": colorVar("surface-bright"),
        "on-secondary-fixed-variant": colorVar("on-secondary-fixed-variant"),
        "surface-tint": colorVar("surface-tint"),
        "on-surface": colorVar("on-surface"),
        "on-primary": colorVar("on-primary"),
        "surface-container": colorVar("surface-container"),
        "on-primary-fixed-variant": colorVar("on-primary-fixed-variant"),
        "on-secondary": colorVar("on-secondary"),
        "inverse-surface": colorVar("inverse-surface"),
        "on-background": colorVar("on-background"),
        "primary-container": colorVar("primary-container"),
        "outline": colorVar("outline"),
        "on-tertiary-fixed": colorVar("on-tertiary-fixed"),
        "surface-dim": colorVar("surface-dim"),
        "secondary": colorVar("secondary"),
        "secondary-fixed": colorVar("secondary-fixed"),
        "error": colorVar("error"),
        "tertiary-fixed-dim": colorVar("tertiary-fixed-dim"),
        "on-tertiary-container": colorVar("on-tertiary-container"),
        "on-error-container": colorVar("on-error-container"),
        "primary-fixed-dim": colorVar("primary-fixed-dim"),
        "inverse-primary": colorVar("inverse-primary"),
        "outline-variant": colorVar("outline-variant"),
        "on-error": colorVar("on-error"),
        "tertiary": colorVar("tertiary"),
        "primary-fixed": colorVar("primary-fixed"),
        "secondary-container": colorVar("secondary-container"),
        "on-tertiary": colorVar("on-tertiary"),
        "on-secondary-container": colorVar("on-secondary-container"),
        "tertiary-container": colorVar("tertiary-container"),
        "surface-container-lowest": colorVar("surface-container-lowest"),
        "secondary-fixed-dim": colorVar("secondary-fixed-dim"),
        "on-surface-variant": colorVar("on-surface-variant"),
        "error-container": colorVar("error-container"),
        "surface-container-low": colorVar("surface-container-low"),
        "surface-container-highest": colorVar("surface-container-highest"),
        "on-secondary-fixed": colorVar("on-secondary-fixed"),
        "on-primary-container": colorVar("on-primary-container"),
        "background": colorVar("background")
      },
      borderRadius: {
        "DEFAULT": "0.5rem",
        "lg": "1rem",
        "xl": "1.5rem",
        "full": "9999px"
      },
      spacing: {
        "2xl": "48px",
        "xl": "32px",
        "lg": "24px",
        "3xl": "64px",
        "md": "16px",
        "xs": "4px",
        "sm": "8px",
        "container-max": "1200px",
        "gutter": "24px"
      },
      fontFamily: {
        "headline": ["Literata", "serif"],
        "display": ["Literata", "serif"],
        "body": ["Nunito Sans", "sans-serif"],
        "label": ["Nunito Sans", "sans-serif"],
        "label-sm": ["Nunito Sans", "sans-serif"],
        "label-md": ["Nunito Sans", "sans-serif"],
        "headline-sm": ["Literata", "serif"],
        "display-lg": ["Literata", "serif"],
        "mono-data": ["Nunito Sans", "sans-serif"],
        "body-lg": ["Nunito Sans", "sans-serif"],
        "body-md": ["Nunito Sans", "sans-serif"],
        "headline-md": ["Literata", "serif"],
        "body-sm": ["Nunito Sans", "sans-serif"],
        "headline-lg": ["Literata", "serif"]
      }
    },
  },
  plugins: [],
}
