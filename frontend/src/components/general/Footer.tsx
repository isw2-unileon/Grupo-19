import React from "react";
import { useNavigate } from "react-router-dom";

// Estructura de datos para que sea facilísimo añadir o quitar enlaces
const footerLinks = [
    {
        title: "ProTracker",
        links: [
            { name: "Sobre nosotros", path: "/about" },
            { name: "Contacto", path: "/contact" },
        ],
    },
    {
        title: "Tus Productos",
        links: [
            { name: "Mis alertas", path: "/notifications" },
            { name: "Guardados", path: "/savedProducts" },
        ],
    },
    {
        title: "Legal",
        links: [
            { name: "Términos de servicio", path: "/terms" },
            { name: "Privacidad", path: "/privacy" },
        ],
    },
];

export default function Footer() {
    const navigate = useNavigate();
    return (
        <footer style={styles.footer}>
            <div style={styles.grid}>
                {footerLinks.map((column, index) => (
                    <div key={index} style={styles.column}>
                        <h4 style={styles.title}>{column.title}</h4>
                        <ul style={styles.list}>
                            {column.links.map((link, linkIndex) => (
                                <li key={linkIndex} style={styles.listItem}>
                                    {/* Si fuera un enlace externo se usaría <a>, aquí usamos el sistema de navegación interno */}
                                    <span
                                        onClick={() => navigate(link.path)}
                                        style={styles.link}
                                    >
                                        {link.name}
                                    </span>
                                </li>
                            ))}
                        </ul>
                    </div>
                ))}
            </div>
            <div style={styles.copyright}>
                © {new Date().getFullYear()} ProTracker. Todos los derechos reservados.
            </div>
        </footer>
    );
}

const styles: Record<string, React.CSSProperties> = {
    footer: {
        backgroundColor: "#1f2937", // Gris muy oscuro
        color: "#f3f4f6", // Texto gris clarito
        padding: "30px 80px 10px 80px", // Espacio interior
        width: "100%",
        boxSizing: "border-box",
        marginTop: "auto", // Ayuda a que se pegue abajo si la página es corta
    },
    grid: {
        display: "grid",
        gridTemplateColumns: "repeat(auto-fit, minmax(200px, 1fr))",
        gap: "40px",
        maxWidth: "1100px",
        margin: "0 auto",
        paddingBottom: "40px",
        borderBottom: "1px solid #374151",
    },
    column: {
        display: "flex",
        flexDirection: "column",
    },
    title: {
        color: "#FACC15", // Toque amarillo para los títulos
        fontSize: "16px",
        fontWeight: "bold",
        marginBottom: "12px",
        textTransform: "uppercase",
        letterSpacing: "0.5px",
    },
    list: {
        listStyle: "none",
        padding: 0,
        margin: 0,
        display: "flex",
        flexDirection: "column",
        gap: "8px",
    },
    listItem: {
        fontSize: "14px",
    },
    link: {
        color: "#9ca3af",
        cursor: "pointer",
        textDecoration: "none",
        transition: "color 0.2s ease",
    },
    copyright: {
        textAlign: "center",
        color: "#6b7280",
        fontSize: "13px",
        marginTop: "10px",
    },
};