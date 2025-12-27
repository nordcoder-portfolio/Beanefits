import type {
    AuthResponse,
    BalanceResponse,
    ClientProfile,
    EventsPage,
    RulesetsPage,
    UsersPage,
    Problem,
} from "@/shared/api/contracts";

import {
    fixtureClientAuth,
    fixtureAdminAuth,
    fixtureClientUser,
    fixtureClientAccount,
    fixtureAdminUser,
    makeUsersPage,
    fixtureCurrentRuleset,
    makeRulesetsPage,
    makeEventsPage,
    makeBalanceFromLatestEvent,
    nextIntId,
    nextUuid,
} from "@/shared/fixtures";

import { ApiError } from "@/shared/api/errors";
import { normalizePhone, PHONE_RE } from "@/shared/lib/phone";

function jsonBody(init: RequestInit): any {
    if (!init.body) return null;
    if (typeof init.body === "string") {
        try { return JSON.parse(init.body); } catch { return null; }
    }
    return null;
}

function problem(status: number, title: string, detail?: string, code?: string, instance?: string): Problem {
    return { type: "about:blank", status, title, detail, code, instance };
}

function requireToken(token?: string) {
    if (!token) throw new ApiError(problem(401, "Unauthorized", "Missing or invalid token"));
}

function requireAdmin(token?: string) {
    requireToken(token);
    if (!String(token).includes(".ADMIN.")) {
        throw new ApiError(problem(403, "Forbidden", "Not enough permissions"));
    }
}

const DEMO_CLIENT_PASSWORD = "client123";
const DEMO_ADMIN_PASSWORD = "admin123";

export async function mockRequestJson<T>(
    path: string,
    init: RequestInit & { token?: string } = {}
): Promise<T> {
    const method = (init.method ?? "GET").toUpperCase();
    const body = jsonBody(init);

    // --- AUTH ---
    if (method === "POST" && path === "/auth/login") {
        const phone = normalizePhone(String(body?.phone ?? ""));
        const password = String(body?.password ?? "");

        if (!PHONE_RE.test(phone)) {
            throw new ApiError(problem(422, "Validation error", "phone: invalid format"));
        }
        if (password.length < 6) {
            throw new ApiError(problem(422, "Validation error", "password: min length is 6"));
        }

        // Демо-аккаунты: строго проверяем креды
        if (phone === fixtureClientUser.phone && password === DEMO_CLIENT_PASSWORD) {
            return fixtureClientAuth as unknown as T;
        }
        if (phone === fixtureAdminUser.phone && password === DEMO_ADMIN_PASSWORD) {
            return fixtureAdminAuth as unknown as T;
        }

        throw new ApiError(problem(401, "Invalid credentials", "phone or password is incorrect"));
    }

    if (method === "POST" && path === "/auth/register") {
        const phone = normalizePhone(String(body?.phone ?? ""));
        const password = String(body?.password ?? "");

        if (!PHONE_RE.test(phone)) {
            throw new ApiError(problem(422, "Validation error", "phone: invalid format"));
        }
        if (password.length < 6) {
            throw new ApiError(problem(422, "Validation error", "password: min length is 6"));
        }
        // условный конфликт (уже есть демо-клиент)
        if (phone === fixtureClientUser.phone || phone === fixtureAdminUser.phone) {
            throw new ApiError(problem(409, "Phone already exists", "phone: already exists"));
        }

        const userId = nextIntId();
        const accountId = nextIntId();
        const publicCode = nextUuid();

        const res: AuthResponse = {
            accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.CLIENT.DEMO",
            user: { ...fixtureClientUser, id: userId, phone, roles: ["CLIENT"], isActive: true },
            account: {
                ...fixtureClientAccount,
                id: accountId,
                publicCode,
                balancePoints: 0,
                totalSpendMoney: "0.00",
                levelCode: "Green Bean",
            },
        };
        return res as unknown as T;
    }

    // --- CLIENT ---
    if (method === "GET" && path === "/me") {
        requireToken(init.token);
        const res: ClientProfile = { user: fixtureClientUser, account: fixtureClientAccount };
        return res as unknown as T;
    }

    if (method === "GET" && path === "/me/balance") {
        requireToken(init.token);
        const page = makeEventsPage(fixtureClientAccount.id, 20);
        const bal: BalanceResponse = makeBalanceFromLatestEvent(fixtureClientAccount.id, page.items);
        return bal as unknown as T;
    }

    if (method === "GET" && path.startsWith("/me/events")) {
        requireToken(init.token);
        const url = new URL(`http://mock${path}`);
        const limit = Number(url.searchParams.get("limit") ?? "20");
        const page: EventsPage = makeEventsPage(fixtureClientAccount.id, limit);
        return page as unknown as T;
    }

    // --- ADMIN ---
    if (method === "GET" && path.startsWith("/admin/users")) {
        requireAdmin(init.token);
        const url = new URL(`http://mock${path}`);
        const limit = Number(url.searchParams.get("limit") ?? "20");
        const page: UsersPage = makeUsersPage(limit);
        return page as unknown as T;
    }

    if (method === "GET" && path === "/admin/rulesets/current") {
        requireAdmin(init.token);
        return fixtureCurrentRuleset as unknown as T;
    }

    if (method === "GET" && path.startsWith("/admin/rulesets")) {
        requireAdmin(init.token);
        const page: RulesetsPage = makeRulesetsPage();
        return page as unknown as T;
    }

    throw new ApiError(problem(404, "Not Found", `No mock handler for ${method} ${path}`));
}
