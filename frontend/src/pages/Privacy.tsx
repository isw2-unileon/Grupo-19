import React from "react";
import Header from "../components/general/Header";
import Footer from "../components/general/Footer";

// Privacy Policy page detailing data handling and cookies
export default function Privacy() {
    return (
        <div style={{ display: "flex", flexDirection: "column", minHeight: "100vh" }}>
            <Header />

            <div style={styles.pageContainer}>
                <div style={styles.card}>
                    <h1 style={styles.title}>Política de Privacidad</h1>

                    <p style={styles.lastUpdated}>Última actualización: {new Date().toLocaleDateString()}</p>

                    <p style={styles.text}>
                        En ProTracker, tu privacidad es de vital importancia para nosotros. Esta Política de Privacidad
                        explica cómo recopilamos, utilizamos y protegemos tu información personal cuando usas nuestra aplicación.
                    </p>

                    <h2 style={styles.sectionTitle}>Información que Recopilamos</h2>
                    <p style={styles.text}>
                        Solo recopilamos la información estrictamente necesaria para proporcionarte nuestros servicios:
                    </p>
                    <ul style={styles.list}>
                        <li style={styles.listItem}><strong>Información de la Cuenta:</strong> Cuando te registras, guardamos tu dirección de correo electrónico y tu nombre de usuario.</li>
                        <li style={styles.listItem}><strong>Datos de Rastreo:</strong> Almacenamos los enlaces de los productos que guardas y las alertas de precio objetivo que configuras.</li>
                    </ul>

                    <h2 style={styles.sectionTitle}>Cómo Utilizamos tu Información</h2>
                    <p style={styles.text}>
                        La información que recopilamos se utiliza exclusivamente para mantener operativo el servicio de ProTracker.
                        Específicamente, utilizamos tu correo para autenticar tu acceso y para enviarte notificaciones cuando un
                        producto baja del precio que has establecido.
                    </p>

                    <h2 style={styles.sectionTitle}>Cookies y Seguridad</h2>
                    <p style={styles.text}>
                        Utilizamos cookies estrictamente necesarias (como tokens JWT) para mantener tu sesión abierta de forma segura.
                        No empleamos cookies de rastreo de terceros con fines publicitarios. Tus contraseñas están fuertemente
                        encriptadas en nuestra base de datos, lo que significa que ni siquiera nuestro equipo puede ver tu contraseña real.
                    </p>

                    <h2 style={styles.sectionTitle}>Intercambio de Datos</h2>
                    <p style={styles.text}>
                        No vendemos, intercambiamos ni alquilamos tu información personal a terceros.
                        Todos tus datos de seguimiento son completamente privados y solo son accesibles desde tu cuenta autenticada.
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
        marginBottom: "8px",
        fontWeight: "bold",
    },
    lastUpdated: {
        fontSize: "14px",
        color: "#9ca3af",
        marginBottom: "24px",
        borderBottom: "2px solid #f3f4f6",
        paddingBottom: "16px",
    },
    sectionTitle: {
        fontSize: "20px",
        color: "#1f2937",
        marginTop: "32px",
        marginBottom: "12px",
        fontWeight: "bold",
    },
    text: {
        fontSize: "16px",
        color: "#4b5563",
        lineHeight: "1.6",
        marginBottom: "16px",
    },
    list: {
        color: "#4b5563",
        fontSize: "16px",
        lineHeight: "1.6",
        marginBottom: "16px",
        paddingLeft: "24px",
    },
    listItem: {
        marginBottom: "8px",
    },
};