import React, { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { Card } from "@/shared/ui/Card";
import { Segmented } from "@/shared/ui/Segmented";
import { Input } from "@/shared/ui/Input";
import { Button } from "@/shared/ui/Button";
import { Callout } from "@/shared/ui/Callout";

import { useAuth } from "@/app/providers/AuthProvider";
import { normalizePhone, PHONE_RE } from "@/shared/lib/phone";
import { humanizeError } from "@/shared/api/errors";

const schema = z.object({
    phone: z.string().refine((v) => PHONE_RE.test(normalizePhone(v)), "Phone: invalid format"),
    password: z.string().min(6, "Password: min 6 chars"),
});

type FormData = z.infer<typeof schema>;

export function AuthPage() {
    const [mode, setMode] = useState<"signin" | "signup">("signin");
    const nav = useNavigate();
    const { login, register } = useAuth();

    const title = mode === "signin" ? "Sign in" : "Sign up";
    const subtitle = "Balance and QR code after login";

    const form = useForm<FormData>({
        resolver: zodResolver(schema),
        defaultValues: { phone: "+7 900 000-00-00", password: "" },
        mode: "onSubmit",
    });

    const isSubmitting = form.formState.isSubmitting;
    const submitError = form.formState.errors.root?.message;

    const segValue = useMemo(() => (mode === "signin" ? "left" : "right") as const, [mode]);

    const onSubmit = form.handleSubmit(async (values) => {
        form.clearErrors("root");

        try {
            const phone = normalizePhone(values.phone);

            const user =
                mode === "signin"
                    ? await login(phone, values.password)
                    : await register(phone, values.password);

            if (user.roles.includes("ADMIN")) {
                nav("/admin/users", { replace: true });
            } else {
                nav("/home", { replace: true });
            }
        } catch (e) {
            form.setError("root", { message: humanizeError(e) });
        }
    });

    return (
        <div className="px-7 pt-14">
            <div className="mb-10">
                <div className="text-4xl font-medium">BeanPoints</div>
                <div className="text-xl text-primary2 mt-2">Coffee loyalty</div>
            </div>

            <Card className="p-8 mb-6">
                <div className="text-4xl font-medium">{title}</div>
                <div className="text-xl text-primary2 mt-2">{subtitle}</div>
            </Card>

            <Segmented
                left="Sign in"
                right="Sign up"
                value={segValue}
                onChange={(v) => setMode(v === "left" ? "signin" : "signup")}
            />

            <form className="mt-10 space-y-8" onSubmit={onSubmit}>
                <div>
                    <div className="text-xl mb-3">Phone</div>
                    <Input {...form.register("phone")} hasError={Boolean(form.formState.errors.phone)} />
                    {form.formState.errors.phone?.message && (
                        <div className="text-danger mt-2">{form.formState.errors.phone.message}</div>
                    )}
                </div>

                <div>
                    <div className="text-xl mb-3">Password</div>
                    <Input
                        type="password"
                        {...form.register("password")}
                        hasError={Boolean(form.formState.errors.password)}
                    />
                    {form.formState.errors.password?.message && (
                        <div className="text-danger mt-2">{form.formState.errors.password.message}</div>
                    )}
                </div>

                {submitError && (
                    <div className="mt-4">
                        <Callout title="Cannot continue">{submitError}</Callout>
                    </div>
                )}

                <Button type="submit" fullWidth disabled={isSubmitting} className="mt-6">
                    Continue
                </Button>
            </form>
        </div>
    );
}
