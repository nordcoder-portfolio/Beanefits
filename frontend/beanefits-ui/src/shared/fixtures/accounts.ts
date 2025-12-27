// src/shared/fixtures/accounts.ts

import type { Account } from "@/shared/api/contracts";
import { isoDaysAgo } from "./clock";
import { nextUuid } from "./ids";

type AccountOverrides = Partial<Account>;

export function makeAccount(overrides: AccountOverrides = {}): Account {
    const id = overrides.id ?? 55;
    return {
        id,
        publicCode: overrides.publicCode ?? nextUuid(),
        balancePoints: overrides.balancePoints ?? 120,
        totalSpendMoney: overrides.totalSpendMoney ?? "12500.50",
        levelCode: overrides.levelCode ?? "Medium Roast",
        createdAt: overrides.createdAt ?? isoDaysAgo(30),
    };
}

export const fixtureClientAccount: Account = makeAccount({
    id: 55,
    publicCode: "550e8400-e29b-41d4-a716-446655440000",
    balancePoints: 240,
    totalSpendMoney: "7800.00",
    levelCode: "Light Roast",
    createdAt: isoDaysAgo(30),
});
