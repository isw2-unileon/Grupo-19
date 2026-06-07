import React from "react";
import Header from "../components/general/Header";
import Footer from "../components/general/Footer";

// Contact page providing communication channels
export default function Contact() {
    return (
        <div style={{ display: "flex", flexDirection: "column", minHeight: "100vh" }}>
            <Header />

            <div style={styles.pageContainer}>
                <div style={styles.card}>
                    <h1 style={styles.title}>Contáctanos</h1>

                    <p style={styles.text}>
                        ¿Tienes alguna pregunta, sugerencia o necesitas ayuda con tu cuenta? Nos encantaría escucharte.
                        Aunque trabajamos constantemente para mejorar ProTracker, la opinión de los usuarios es nuestro activo más valioso.
                    </p>

                    <div style={styles.contactInfoBox}>
                        <h2 style={styles.sectionTitle}>Ponte en contacto</h2>
                        <p style={styles.text}>
                            Para consultas generales, soporte técnico o sugerencias de nuevas funcionalidades, por favor envíanos un correo electrónico.
                            Nuestro objetivo es responder a todas las consultas en un plazo de 24 a 48 horas laborables.
                        </p>

                        <div style={styles.emailWrapper}>
                            <strong>Soporte por Correo: </strong>
                            {/* Standard mailto link to open the user's default email client */}
                            <a href="mailto:support@protracker.com" style={styles.link}>
                                support@protracker.com
                            </a>
                        </div>
                    </div>

                    <p style={styles.text}>
                        <em>Nota: Asegúrate de incluir la dirección de correo electrónico de tu cuenta en el mensaje si nos contactas por un problema específico asociado a tu perfil.</em>
                    </p>
                </div>
            </div>

            <Footer />
        </div>
    );
}

const styles: Record<string, React.CSSProperties> = {
    pageContainer: {
        padding: "40px 20px",
        display: "flex",
        justifyContent: "center",
        backgroundColor: "#f9fafb",
        flexGrow: 1,
    },
    card: {
        background: "white",
        borderRadius: "12px",
        padding: "40px",
        maxWidth: "800px",
        width: "100%",
        boxShadow: "0 4px 15px rgba(0,0,0,0.02)",
        border: "1px solid #eee",
        borderTop: "4px solid #FACC15",
    },
    title: {
        fontSize: "28px",
        color: "#1f2937",
        marginBottom: "24px",
        fontWeight: "bold",
        borderBottom: "2px solid #f3f4f6",
        paddingBottom: "16px",
    },
    sectionTitle: {
        fontSize: "20px",
        color: "#1f2937",
        marginTop: "10px",
        marginBottom: "12px",
        fontWeight: "bold",
    },
    text: {
        fontSize: "16px",
        color: "#4b5563",
        lineHeight: "1.6",
        marginBottom: "16px",
    },
    contactInfoBox: {
        backgroundColor: "#f9fafb",
        padding: "20px",
        borderRadius: "8px",
        border: "1px solid #e5e7eb",
        margin: "30px 0",
    },
    emailWrapper: {
        fontSize: "16px",
        color: "#1f2937",
        marginTop: "16px",
    },
    link: {
        color: "#2563eb",
        textDecoration: "underline",
        fontWeight: "bold",
    },
};