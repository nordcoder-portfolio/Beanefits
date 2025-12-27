// src/shared/fixtures/clock.ts

export const FIXTURE_NOW_ISO = "2025-12-24T12:00:00.000Z" as const;

export function isoMinutesAgo(minutes: number): string {
    const base = new Date(FIXTURE_NOW_ISO).getTime();
    return new Date(base - minutes * 60_000).toISOString();
}

export function isoDaysAgo(days: number): string {
    const base = new Date(FIXTURE_NOW_ISO).getTime();
    return new Date(base - days * 24 * 60 * 60_000).toISOString();
}
