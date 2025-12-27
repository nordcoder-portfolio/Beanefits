import React, { useEffect, useMemo, useState } from "react";
import { Outlet } from "react-router-dom";

const DESIGN_W = 430; // ширина макета "телефона"
const DESIGN_H = 844; // высота макета "телефона" (iPhone-like)

function clamp(n: number, min: number, max: number) {
    return Math.max(min, Math.min(max, n));
}

export function PhoneFrameLayout({
                                     withNav = false,
                                     nav,
                                 }: {
    withNav?: boolean;
    nav?: React.ReactNode;
}) {
    const [scale, setScale] = useState(1);

    useEffect(() => {
        const recalc = () => {
            // px-4 (16*2) + py-8 (32*2) как в layout ниже
            const availableW = window.innerWidth - 32;
            const availableH = window.innerHeight - 64;

            const sW = availableW / DESIGN_W;
            const sH = availableH / DESIGN_H;

            // не увеличиваем больше 1, только уменьшаем при необходимости
            const next = clamp(Math.min(sW, sH, 1), 0.5, 1);
            setScale(next);
        };

        recalc();
        window.addEventListener("resize", recalc);
        return () => window.removeEventListener("resize", recalc);
    }, []);

    const outerSize = useMemo(
        () => ({
            width: Math.round(DESIGN_W * scale),
            height: Math.round(DESIGN_H * scale),
        }),
        [scale]
    );

    return (
        <div className="h-dvh w-full bg-bg flex items-center justify-center px-4 py-8 overflow-hidden">
            {/* контейнер фиксированного "видимого" размера */}
            <div style={outerSize} className="relative">
                {/* реальный макет в дизайн-размере, масштабируем через transform */}
                <div
                    style={{
                        width: DESIGN_W,
                        height: DESIGN_H,
                        transform: `scale(${scale})`,
                        transformOrigin: "top left",
                    }}
                    className="rounded-[44px] bg-bg shadow-soft overflow-hidden relative"
                >
                    <div className={(withNav ? "pb-24 " : "") + "h-full overflow-y-auto"}>
                        <Outlet/>
                    </div>

                    {withNav && (
                        <div className="absolute inset-x-0 bottom-0 bg-card border-t border-border">
                            {nav}
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
