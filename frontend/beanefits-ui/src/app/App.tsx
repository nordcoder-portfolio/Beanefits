import React from "react";
import {createBrowserRouter, Navigate, RouterProvider} from "react-router-dom";
import {QueryProvider} from "@/app/providers/QueryProvider";
import {AuthProvider} from "@/app/providers/AuthProvider";
import {RequireAuth} from "@/app/guards/RequireAuth";
import {PhoneFrameLayout} from "@/widgets/layout/PhoneFrameLayout";
import {BottomNav} from "@/widgets/nav/BottomNav";

import {AuthPage} from "@/pages/auth/AuthPage";
import {HomePage} from "@/pages/client/HomePage";
import {HistoryPage} from "@/pages/client/HistoryPage";
import {ProfilePage} from "@/pages/client/ProfilePage";
import { AdminLayout } from "@/widgets/layout/AdminLayout";
import { AdminUsersPage } from "@/pages/admin/AdminUsersPage";
import { AdminRulesPage } from "@/pages/admin/AdminRulesPage";


const router = createBrowserRouter([
    {
        element: <PhoneFrameLayout/>,
        children: [
            {index: true, element: <Navigate to="/auth" replace/>},
            {path: "/auth", element: <AuthPage/>},
        ],
    },
    {
        element: (
            <RequireAuth roles={["CLIENT"]}>
                <PhoneFrameLayout withNav nav={<BottomNav/>}/>
            </RequireAuth>
        ),
        children: [
            {path: "/home", element: <HomePage/>},
            {path: "/history", element: <HistoryPage/>},
            {path: "/profile", element: <ProfilePage/>},
            {path: "*", element: <Navigate to="/home" replace/>},
        ],
    },
    {
        path: "/admin",
        element: (
            <RequireAuth roles={["ADMIN"]}>
                <AdminLayout />
            </RequireAuth>
        ),
        children: [
            { index: true, element: <Navigate to="/admin/users" replace /> },
            { path: "users", element: <AdminUsersPage /> },
            { path: "rules", element: <AdminRulesPage /> },
        ],
    },
]);

export function App() {
    return (
        <QueryProvider>
            <AuthProvider>
                <RouterProvider router={router}/>
            </AuthProvider>
        </QueryProvider>
    );
}
