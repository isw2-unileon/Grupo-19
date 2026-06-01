import React from "react";
import { useNavigate } from "react-router-dom";
// Asegúrate de que la ruta al logo es correcta según donde guardes este componente
import logoImage from "../../assets/Logo.png";

export default function Header() {
    const iconStyle = { width: "24px", height: "24px", cursor: "pointer", color: "#1f2937" };
    const navigate = useNavigate();

    return (
        <header style={styles.header}>
            {/* Zona Izquierda: Logo */}
            <div style={styles.logoContainer} onClick={() => navigate("/mainPage")}>
                <img src={logoImage} alt="ProTracker Logo" style={styles.logoImg} />
            </div>

            {/* Zona Derecha: Acciones */}
            <div style={{ display: "flex", gap: "24px", alignItems: "center" }}>

                {/* Botón: Notificaciones (Campana) */}
                <svg
                    onClick={() => navigate("/notifications")}
                    style={iconStyle} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}
                >
                    <path strokeLinecap="round" strokeLinejoin="round" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                </svg>

                {/* Botón: Guardados (Marcador/Bookmark) - SUSTITUYE AL SOBRE */}
                <svg
                    onClick={() => navigate("/savedProducts")}
                    style={iconStyle} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}
                >
                    <path strokeLinecap="round" strokeLinejoin="round" d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z" />
                </svg>

                {/* Botón: Perfil (Usuario) */}
                <div
                    onClick={() => navigate("/profile")}
                    style={{ padding: "8px", borderRadius: "50%", background: "rgba(255,255,255,0.4)", display: "flex", cursor: "pointer" }}
                >
                    <svg style={iconStyle} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                        <path strokeLinecap="round" strokeLinejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                </div>

            </div>
        </header>
    );
}

const styles: Record<string, React.CSSProperties> = {
    header: {
        display: "flex",
        justifyContent: "space-between",
        alignItems: "center",
        backgroundColor: "#FACC15", // Amarillo corporativo
        padding: "10px 80px", // Mismos márgenes laterales que tu MainLayout
        width: "100%",
        boxSizing: "border-box",
        boxShadow: "0 2px 10px rgba(0,0,0,0.05)",
    },
    logoContainer: {
        display: "flex",
        alignItems: "center",
        cursor: "pointer",
    },
    logoImg: {
        height: "60px", // Ligeramente más pequeño para que encaje bien en la barra superior
        width: "auto",
        objectFit: "contain",
    },
};