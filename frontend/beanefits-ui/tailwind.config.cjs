/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./index.html", "./src/**/*.{ts,tsx}"],
    theme: {
        extend: {
            colors: {
                bg: "rgb(var(--bp-bg) / <alpha-value>)",
                card: "rgb(var(--bp-card) / <alpha-value>)",
                ink: "rgb(var(--bp-ink) / <alpha-value>)",
                muted: "rgb(var(--bp-muted) / <alpha-value>)",
                primary: "rgb(var(--bp-primary) / <alpha-value>)",
                primary2: "rgb(var(--bp-primary-2) / <alpha-value>)",
                border: "rgb(var(--bp-border) / <alpha-value>)",
                accent: "rgb(var(--bp-accent) / <alpha-value>)",
                peach: "rgb(var(--bp-peach) / <alpha-value>)",
                success: "rgb(var(--bp-success) / <alpha-value>)",
                danger: "rgb(var(--bp-danger) / <alpha-value>)",
            },
            borderRadius: {
                xl2: "28px",
                xl3: "36px",
            },
            boxShadow: {
                soft: "0 14px 50px rgba(0,0,0,0.08)",
            },
        },
    },
    plugins: [],
};
