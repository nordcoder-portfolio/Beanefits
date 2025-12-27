// src/shared/fixtures/users.ts

import type { User, UsersPage } from "@/shared/api/contracts";
import { isoDaysAgo } from "./clock";
import { nextIntId } from "./ids";

type UserOverrides = Partial<User>;

export function makeUser(overrides: UserOverrides = {}): User {
    const id = overrides.id ?? nextIntId();
    return {
        id,
        phone: overrides.phone ?? `+7900${String(id).padStart(7, "0")}`,
        roles: overrides.roles ?? ["CLIENT"],
        isActive: overrides.isActive ?? true,
        createdAt: overrides.createdAt ?? isoDaysAgo(10),
    };
}

export const fixtureClientUser: User = makeUser({
    id: 10,
    phone: "+79000000000",
    roles: ["CLIENT"],
    createdAt: isoDaysAgo(30),
});

export const fixtureAdminUser: User = makeUser({
    id: 1,
    phone: "+79990000000",
    roles: ["ADMIN"],
    createdAt: isoDaysAgo(120),
});

export const fixtureCashierUser: User = makeUser({
    id: 2,
    phone: "+79990000001",
    roles: ["CASHIER"],
    createdAt: isoDaysAgo(90),
});

export function makeUsersPage(count = 20): UsersPage {
    const items: User[] = [
        fixtureAdminUser,
        fixtureCashierUser,
        fixtureClientUser,
        ...Array.from({ length: Math.max(0, count - 3) }, () => makeUser()),
    ];

    return { items, total: items.length };
}
