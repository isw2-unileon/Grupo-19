import React from "react";
import { useNavigate } from "react-router-dom";

// Data structure to make it super easy to add or remove links
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
                                    {/* If it were an external link, <a> would be used; here we use the internal navigation system */}
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
        backgroundColor: "#1f2937",
        color: "#f3f4f6", 
        padding: "30px 80px 10px 80px", // Interior space
        width: "100%",
        boxSizing: "border-box",
        marginTop: "auto", // Helps it stick at the bottom if the page is short
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
        color: "#FACC15", 
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