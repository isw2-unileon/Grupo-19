import React from "react";
import Header from "../components/general/Header";
import Footer from "../components/general/Footer";

// About page component detailing the application's purpose
export default function About() {
    return (
        <div style={{ display: "flex", flexDirection: "column", minHeight: "100vh" }}>
            <Header />

            <div style={styles.pageContainer}>
                <div style={styles.card}>
                    <h1 style={styles.title}>Sobre ProTracker</h1>

                    <p style={styles.text}>
                        Bienvenido a ProTracker, tu herramienta definitiva para compras inteligentes online.
                        Nuestra misión es ayudarte a ahorrar tiempo y dinero automatizando el proceso
                        de rastreo de precios de productos en diversas tiendas online.
                    </p>

                    <h2 style={styles.sectionTitle}>Nuestra Visión</h2>
                    <p style={styles.text}>
                        Creemos que todo el mundo merece conseguir las mejores ofertas sin tener que
                        revisar manualmente las páginas web todos los días. ProTracker hace el trabajo pesado por ti:
                        simplemente añade el enlace de un producto y nosotros monitorizaremos las bajadas de precio,
                        notificándote en el momento en que alcance tu precio objetivo.
                    </p>

                    <h2 style={styles.sectionTitle}>El Equipo</h2>
                    <p style={styles.text}>
                        Desarrollado por un dedicado equipo de estudiantes de ingeniería informática, ProTracker
                        fue construido teniendo en mente el rendimiento, la seguridad y la experiencia del usuario.
                        Utilizamos tecnologías web modernas para asegurar que tus datos estén a salvo y que nuestro motor
                        de extracción de datos sea lo más eficiente posible.
                    </p>
                </div>
            </div>

            <Footer />
        </div>
    );
}

// Reusable styles matching the application theme
const styles: Record<string, React.CSSProperties> = {
    pageContainer: {
        padding: "40px 20px",
        display: "flex",
        justifyContent: "center",
        backgroundColor: "#f9fafb", // Light background to make the card pop
        flexGrow: 1, // Pushes the footer to the bottom of the screen
    },
    card: {
        background: "white",
        borderRadius: "12px",
        padding: "40px",
        maxWidth: "800px",
        width: "100%",
        boxShadow: "0 4px 15px rgba(0,0,0,0.02)",
        border: "1px solid #eee",
        borderTop: "4px solid #FACC15", // Signature yellow accent
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
};