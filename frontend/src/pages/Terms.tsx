import React from "react";
import Header from "../components/general/Header";
import Footer from "../components/general/Footer";

// Terms of Service page with standard boilerplate text
export default function Terms() {
    return (
        <div style={{ display: "flex", flexDirection: "column", minHeight: "100vh" }}>
            <Header />

            <div style={styles.pageContainer}>
                <div style={styles.card}>
                    <h1 style={styles.title}>Términos de Servicio</h1>

                    <p style={styles.lastUpdated}>Última actualización: {new Date().toLocaleDateString()}</p>

                    <p style={styles.text}>
                        Por favor, lee estos Términos de Servicio cuidadosamente antes de utilizar la página web
                        y el servicio de ProTracker operado por nuestro equipo.
                    </p>

                    <h2 style={styles.sectionTitle}>1. Aceptación de los Términos</h2>
                    <p style={styles.text}>
                        Al acceder o utilizar nuestro servicio, aceptas estar sujeto a estos Términos.
                        Si no estás de acuerdo con alguna parte de los términos, no podrás acceder a la plataforma.
                    </p>

                    <h2 style={styles.sectionTitle}>2. Descripción del Servicio</h2>
                    <p style={styles.text}>
                        ProTracker proporciona una herramienta web que permite a los usuarios rastrear los precios de
                        productos en tiendas online de terceros. No vendemos productos directamente, ni podemos garantizar
                        la exactitud de los precios en el momento exacto de la compra, ya que las tiendas de origen actualizan
                        sus datos de forma constante e independiente.
                    </p>

                    <h2 style={styles.sectionTitle}>3. Cuentas de Usuario</h2>
                    <p style={styles.text}>
                        Cuando creas una cuenta con nosotros, debes proporcionar información que sea precisa, completa
                        y actual en todo momento. No hacerlo constituye un incumplimiento de los Términos, lo que puede resultar
                        en la suspensión inmediata de tu cuenta en nuestro servicio.
                    </p>

                    <h2 style={styles.sectionTitle}>4. Uso Aceptable</h2>
                    <p style={styles.text}>
                        Aceptas no utilizar el servicio para abusar, acosar o interferir con nuestros servidores o red.
                        La extracción de datos masiva y automatizada de nuestra propia plataforma, o el intento de eludir nuestras
                        medidas de seguridad, está estrictamente prohibido.
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
};