//import React from "react";
import Header from "../components/general/Header";
import Footer from "../components/general/Footer";

export default function Profile() {
    return (
        <div style={{ display: "flex", flexDirection: "column", minHeight: "100vh", backgroundColor: "#fafafa" }}>
            <Header />

            <main style={{ flex: 1, display: "flex", justifyContent: "center", alignItems: "center" }}>
                <h2>Perfil</h2>
            </main>

            <Footer />
        </div>
    );
}