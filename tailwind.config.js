/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
      'templates/*.templ',
  ],
  theme: {
    extend: {
      container: {
        center: true,
        padding: {
            DEFAULT: '12px',
        },
        screens: {
          sm: '680px',
          md: '680px',
          lg: '744px',
          xl: '744px',
        },
      },
      colors: {
        background: {
          DEFAULT: "oklch(var(--background) / <alpha-value>)",
        },
        foreground: {
          DEFAULT: "oklch(var(--foreground) / <alpha-value>)",
        },
        border: {
          DEFAULT: "oklch(var(--border) / <alpha-value>)",
        },
        alt: {
          DEFAULT: "oklch(var(--alt-background) / <alpha-value>)",
          foreground: "oklch(var(--alt-foreground) / <alpha-value>)",
        },
        sidebar: {
          DEFAULT: "oklch(var(--sidebar-background) / <alpha-value>)",
          foreground: "oklch(var(--sidebar-foreground) / <alpha-value>)",
        },
        primary: {
          DEFAULT: "oklch(var(--primary) / <alpha-value>)",
        },
        secondary: {
          DEFAULT: "oklch(var(--secondary-background) / <alpha-value>)",
          foreground: "oklch(var(--secondary-foreground) / <alpha-value>)",
        },
        success: {
          DEFAULT: "oklch(var(--success) / <alpha-value>)",
        },
        error: {
          DEFAULT: "oklch(var(--error) / <alpha-value>)",
        },
        warning: {
          DEFAULT: "oklch(var(--warning) / <alpha-value>)",
        },
        black: {
          DEFAULT: "oklch(var(--black) / <alpha-value>)",
        },
      },
      fontFamily: {
        sans: ["var(--font-sans)"],
        mono: ["var(--font-mono)"],
      },
      fontSize: {
        xss: ["7px", { lineHeight: "10px", letterSpacing: "0.005em" }],
        xs: ["9px", { lineHeight: "11px", letterSpacing: "0.005em" }],
        sm: ["11px", { lineHeight: "13px", letterSpacing: "0.005em" }],
        md: ["12px", { lineHeight: "15px", letterSpacing: "0.0025em" }],
        base: ["14px", { lineHeight: "17px", letterSpacing: "0em" }],
        lg: ["16px", { lineHeight: "24px", letterSpacing: "0em" }],
        xl: ["18px", { lineHeight: "26px", letterSpacing: "-0.0025em" }],
        "2xl": ["20px", { lineHeight: "28px", letterSpacing: "-0.005em" }],
        "3xl": ["24px", { lineHeight: "28px", letterSpacing: "-0.00625em" }],
        "4xl": ["28px", { lineHeight: "36px", letterSpacing: "-0.0075em" }],
        "5xl": ["35px", { lineHeight: "40px", letterSpacing: "-0.01em" }],
        "6xl": ["45px", { lineHeight: "48px", letterSpacing: "-0.0125em" }],
        "7xl": ["60px", { lineHeight: "60px", letterSpacing: "-0.025em" }],
      },
      fontWeight: {
        light: "300",
        normal: "400",
        medium: "400",
        semibold: "600",
        bold: "700",
      },
      borderRadius: {
        DEFAULT: "var(--border-radius)",
      },
      dropShadow: {
        DEFAULT: "rgba(66, 66, 90, 0.6) 0px 0px 0px 1px, rgba(11, 14, 22, 0.15) 0px 4px 10px 0px, rgba(11, 14, 22, 0.25) 0px 8px 20px 0px;",
      },
    },
  },
  plugins: [],
  safelist: [
        'toast',
        'toast-success',
        'toast-error',
        'toast-warning',
        'markdown-alert',
        'markdown-alert-title',
        'bg-black',
        'message-edited',
        'message-box-container',
    {
      pattern: /markdown-alert/,
    },
    {
      pattern: /markdown-alert-(tip|note|important|warning|caution)/,
    },
    {
      pattern: /hl-(chroma|line|kd|p|c|o|nf|kt|k)/,
    },
  ]
}

