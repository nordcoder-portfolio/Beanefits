import React from "react";
import { Card } from "@/shared/ui/Card";
import { Button } from "@/shared/ui/Button";
import { useAuth } from "@/app/providers/AuthProvider";
import { useMeQuery } from "@/entities/me/queries";

export function ProfilePage() {
    const { logout } = useAuth();
    const me = useMeQuery();

    const phone = me.data?.user.phone ?? "+7 900 000-00-00";

    return (
        <div className="px-7 pt-10 space-y-8">
            <div>
                <div className="text-4xl font-medium">Profile</div>
                <div className="text-xl text-muted mt-2">Account settings</div>
            </div>

            <Card className="p-8">
                <div className="text-xl text-muted">Phone</div>
                <div className="text-2xl mt-2">{phone}</div>
            </Card>

            <Card className="p-8">
                <div className="text-xl text-muted">Actions</div>
                <Button fullWidth className="mt-6 bg-primary text-white" type="button" onClick={logout}>
                    Logout
                </Button>
            </Card>
        </div>
    );
}
