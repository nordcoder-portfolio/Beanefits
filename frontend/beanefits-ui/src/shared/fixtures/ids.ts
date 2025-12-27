// src/shared/fixtures/ids.ts

let intId = 1000;
export function nextIntId(): number {
    intId += 1;
    return intId;
}

/**
 * Детерминированные UUID для фикстур (не используем random в тестах).
 * При необходимости — расширим пул.
 */
const UUIDS = [
    "550e8400-e29b-41d4-a716-446655440000",
    "550e8400-e29b-41d4-a716-446655440001",
    "550e8400-e29b-41d4-a716-446655440002",
    "550e8400-e29b-41d4-a716-446655440003",
    "550e8400-e29b-41d4-a716-446655440004",
    "550e8400-e29b-41d4-a716-446655440005",
    "550e8400-e29b-41d4-a716-446655440006",
    "550e8400-e29b-41d4-a716-446655440007",
] as const;

let uuidIdx = 0;
export function nextUuid(): string {
    const u = UUIDS[uuidIdx % UUIDS.length];
    uuidIdx += 1;
    return u;
}
