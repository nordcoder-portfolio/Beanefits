export function normalizePhone(input: string): string {
    let s = input.trim().replace(/[^\d+]/g, "");
    if (s.includes("+") && !s.startsWith("+")) s = s.replace(/\+/g, "");
    return s;
}

export const PHONE_RE = /^\+?[1-9]\d{10,14}$/;
