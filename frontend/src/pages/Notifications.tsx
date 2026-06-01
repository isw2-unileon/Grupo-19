import React from "react";
import Header from "../components/general/Header";
import Footer from "../components/general/Footer";

export default function Notifications() {
    return (
        <div style={{ display: "flex", flexDirection: "column", minHeight: "100vh", backgroundColor: "#fafafa" }}>
            <Header />

            <main style={{ flex: 1, display: "flex", justifyContent: "center", alignItems: "center" }}>
                <h2>Notificaciones</h2>
            </main>

            <Footer />
        </div>
    );
}