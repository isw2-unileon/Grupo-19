import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import logoImage from "../../assets/Logo.png";

export default function Header() {
    const iconStyle = { width: "24px", height: "24px", cursor: "pointer", color: "#1f2937" };
    const navigate = useNavigate();

    // Estado para controlar la visibilidad del modal
    const [showLogoutModal, setShowLogoutModal] = useState(false);

    // Función para manejar el cierre de sesión real
    const handleLogout = async () => {
        try {
            const response = await fetch("http://localhost:8080/api/logout", {
                method: "POST",
                credentials: "include",
            });

            if (response.ok) {
                window.location.href = "/";
            } else {
                console.error("Error al intentar cerrar sesión");
                setShowLogoutModal(false); // Ocultamos el modal si falla
            }
        } catch (error) {
            console.error("Error de red al cerrar sesión:", error);
            setShowLogoutModal(false); // Ocultamos el modal si falla
        }
    };

    return (
        <>
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

                    {/* Botón: Guardados (Marcador/Bookmark) */}
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
                        title="Mi Perfil"
                    >
                        <svg style={iconStyle} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                        </svg>
                    </div>

                    {/* Botón: Cerrar Sesión (Abre el modal) */}
                    <div
                        onClick={() => setShowLogoutModal(true)}
                        style={{
                            padding: "8px",
                            borderRadius: "50%",
                            background: "#ffffff",
                            display: "flex",
                            cursor: "pointer",
                            boxShadow: "0 1px 3px rgba(0,0,0,0.1)"
                        }}
                        title="Cerrar sesión"
                    >
                        <svg style={iconStyle} fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                        </svg>
                    </div>

                </div>
            </header>

            {/* Modal de confirmación */}
            {showLogoutModal && (
                <div style={styles.modalOverlay}>
                    <div style={styles.modalContent}>
                        <h3 style={{ marginTop: 0, color: "#1f2937" }}>Cerrar sesión</h3>
                        <p style={{ color: "#4b5563", marginBottom: "24px" }}>
                            ¿Estás seguro de que quieres salir de tu cuenta?
                        </p>
                        <div style={styles.modalActions}>
                            <button
                                onClick={() => setShowLogoutModal(false)}
                                style={styles.cancelButton}
                            >
                                Cancelar
                            </button>
                            <button
                                onClick={handleLogout}
                                style={styles.confirmButton}
                            >
                                Sí, salir
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </>
    );
}

const styles: Record<string, React.CSSProperties> = {
    header: {
        display: "flex",
        justifyContent: "space-between",
        alignItems: "center",
        backgroundColor: "#FACC15",
        padding: "10px 80px",
        width: "100%",
        boxSizing: "border-box",
        boxShadow: "0 2px 10px rgba(0,0,0,0.05)",
        height: "80px",
    },
    logoContainer: {
        display: "flex",
        alignItems: "center",
        cursor: "pointer",
    },
    logoImg: {
        height: "105px",
        width: "auto",
        objectFit: "contain",
        marginTop: "5px",
    },
    // --- ESTILOS DEL MODAL ---
    modalOverlay: {
        position: "fixed",
        top: 0,
        left: 0,
        right: 0,
        bottom: 0,
        backgroundColor: "rgba(0, 0, 0, 0.5)", // Fondo oscuro semitransparente
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        zIndex: 9999, // Asegura que se dibuje por encima de todo
    },
    modalContent: {
        backgroundColor: "#ffffff",
        padding: "24px",
        borderRadius: "8px",
        boxShadow: "0 4px 6px rgba(0, 0, 0, 0.1)",
        maxWidth: "400px",
        width: "90%",
        textAlign: "center",
        fontFamily: "sans-serif",
    },
    modalActions: {
        display: "flex",
        justifyContent: "center",
        gap: "16px",
    },
    cancelButton: {
        padding: "10px 20px",
        border: "1px solid #d1d5db",
        backgroundColor: "#ffffff",
        color: "#374151",
        borderRadius: "6px",
        cursor: "pointer",
        fontWeight: "bold",
    },
    confirmButton: {
        padding: "10px 28px",
        border: "none",
        backgroundColor: "#FACC15",
        color: "#1f2937",
        borderRadius: "6px",
        cursor: "pointer",
        fontWeight: "bold",
    },
};