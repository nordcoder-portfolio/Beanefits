export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080";

export const USE_FIXTURES = import.meta.env.VITE_USE_FIXTURES === "1";

export const API_TIMEOUT_MS = Number(import.meta.env.VITE_API_TIMEOUT_MS ?? "10000");

export const API_DEBUG = import.meta.env.DEV && import.meta.env.VITE_API_DEBUG === "1";
