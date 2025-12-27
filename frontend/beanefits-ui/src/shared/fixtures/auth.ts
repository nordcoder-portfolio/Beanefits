// src/shared/fixtures/auth.ts

import type { AuthResponse } from "@/shared/api/contracts";
import { fixtureClientUser, fixtureAdminUser } from "./users";
import { fixtureClientAccount } from "./accounts";

const DEMO_JWT_CLIENT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.CLIENT.DEMO";
const DEMO_JWT_ADMIN  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.ADMIN.DEMO";

export const fixtureClientAuth: AuthResponse = {
    accessToken: DEMO_JWT_CLIENT,
    user: fixtureClientUser,
    account: fixtureClientAccount,
};

export const fixtureAdminAuth: AuthResponse = {
    accessToken: DEMO_JWT_ADMIN,
    user: fixtureAdminUser,
    account: {
        ...fixtureClientAccount,
        id: 1,
        balancePoints: 0,
        totalSpendMoney: "0.00",
        levelCode: "Green Bean",
    },
};
