import React from "react";
import { Navigate, useLocation } from "react-router-dom";
import type { RoleCode } from "@/shared/api/contracts";
import { useAuth } from "@/app/providers/AuthProvider";

export function RequireAuth({ children, roles }: { children: React.ReactNode; roles?: RoleCode[] }) {
    const { accessToken, user } = useAuth();
    const loc = useLocation();

    if (!accessToken || !user) return <Navigate to="/auth" replace state={{ from: loc.pathname }} />;
    if (roles && !roles.some((r) => user.roles.includes(r))) return <Navigate to="/auth" replace />;

    return <>{children}</>;
}
