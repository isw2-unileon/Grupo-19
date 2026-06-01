import React, { useState } from "react";
import Header from "../components/general/Header";
import Footer from "../components/general/Footer";

export default function Profile() {
  // --- Estados para los campos del Perfil ---
  const [username, setUsername] = useState("Usuario_Grupo19");
  const [email, setEmail] = useState("grupo19@unileon.es");

  // --- Estados para el cambio de Contraseña ---
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  // --- Funciones preparadas para conectar con el Backend en Go ---
  const handleUpdateProfile = (e: React.FormEvent) => {
    e.preventDefault();
    console.log("Enviando a Go los nuevos datos de usuario:", { username, email });
    alert("Perfil actualizado correctamente (Simulado)");
  };

  const handleUpdatePassword = (e: React.FormEvent) => {
    e.preventDefault();
    if (newPassword !== confirmPassword) {
      alert("Las contraseñas nuevas no coinciden");
      return;
    }
    console.log("Enviando a Go el cambio de contraseña");
    alert("Contraseña actualizada correctamente (Simulado)");
    // Limpiar campos de seguridad
    setCurrentPassword("");
    setNewPassword("");
    setConfirmPassword("");
  };

  return (
    <div style={{ display: "flex", flexDirection: "column", minHeight: "100vh", backgroundColor: "#fafafa", fontFamily: "sans-serif" }}>
      {/*  cabecera global  */}
      <Header />

      {/* Contenido principal  */}
      <main style={styles.mainContent}>
        <h2 style={styles.pageTitle}>Ajustes de Cuenta</h2>

        <div style={styles.sectionsGrid}>
          {/* INFORMACIÓN DEL PERFIL */}
          <section style={styles.card}>
            <h3 style={styles.cardTitle}>Información del Perfil</h3>
            <form onSubmit={handleUpdateProfile} style={styles.form}>
              <div style={styles.inputGroup}>
                <label style={styles.label}>Nombre de Usuario</label>
                <input
                  type="text"
                  value={username}
                  onChange={(e) => setUsername(e.target.value)}
                  style={styles.input}
                  required
                />
              </div>

              <div style={styles.inputGroup}>
                <label style={styles.label}>Correo Electrónico</label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  style={styles.input}
                  required
                />
              </div>

              <button type="submit" style={styles.primaryButton}>
                Guardar Cambios
              </button>
            </form>
          </section>

          {/* SEGURIDAD / CONTRASEÑA */}
          <section style={styles.card}>
            <h3 style={styles.cardTitle}>Seguridad y Credenciales</h3>
            <form onSubmit={handleUpdatePassword} style={styles.form}>
              <div style={styles.inputGroup}>
                <label style={styles.label}>Contraseña Actual</label>
                <input
                  type="password"
                  placeholder="••••••••"
                  value={currentPassword}
                  onChange={(e) => setCurrentPassword(e.target.value)}
                  style={styles.input}
                  required
                />
              </div>

              <div style={styles.inputGroup}>
                <label style={styles.label}>Nueva Contraseña</label>
                <input
                  type="password"
                  placeholder="Mínimo 6 caracteres"
                  value={newPassword}
                  onChange={(e) => setNewPassword(e.target.value)}
                  style={styles.input}
                  required
                />
              </div>

              <div style={styles.inputGroup}>
                <label style={styles.label}>Confirmar Nueva Contraseña</label>
                <input
                  type="password"
                  placeholder="Repite tu nueva contraseña"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  style={styles.input}
                  required
                />
              </div>

              <button type="submit" style={styles.primaryButton}>
                Actualizar Contraseña
              </button>
            </form>
          </section>
        </div>
      </main>

      {/* pie de página global */}
      <Footer />
    </div>
  );
}


const styles: Record<string, React.CSSProperties> = {
  mainContent: {
    flex: 1,
    width: "100%",
    maxWidth: "1100px",
    margin: "0 auto",
    padding: "40px 40px 80px 40px",
    boxSizing: "border-box",
  },
  pageTitle: {
    fontSize: "26px",
    color: "#1f2937",
    marginBottom: "32px",
    fontWeight: "bold",
    borderLeft: "5px solid #FACC15", // Línea vertical amarilla identitaria
    paddingLeft: "12px",
    lineHeight: "1",
  },
  sectionsGrid: {
    display: "grid",
    gridTemplateColumns: "repeat(auto-fit, minmax(450px, 1fr))", // Responsivo automático
    gap: "32px",
  },
  card: {
    background: "white",
    borderRadius: "12px",
    padding: "32px",
    boxShadow: "0 4px 15px rgba(0,0,0,0.01)",
    border: "1px solid #eee",
    borderTop: "4px solid #FACC15", // Detalle superior amarillo corporativo
  },
  cardTitle: {
    margin: "0 0 24px 0",
    fontSize: "18px",
    color: "#1f2937",
    fontWeight: "bold",
  },
  form: {
    display: "flex",
    flexDirection: "column",
    gap: "20px",
  },
  inputGroup: {
    display: "flex",
    flexDirection: "column",
    gap: "8px",
  },
  label: {
    fontSize: "14px",
    fontWeight: "600",
    color: "#4b5563",
  },
  input: {
    padding: "12px 16px",
    fontSize: "14px",
    border: "1px solid #ddd",
    borderRadius: "8px",
    outline: "none",
    backgroundColor: "#fff",
  },
  primaryButton: {
    padding: "14px 20px",
    border: "none",
    background: "#FACC15",
    color: "black",
    fontWeight: "bold",
    fontSize: "14px",
    borderRadius: "8px",
    cursor: "pointer",
    marginTop: "10px",
    alignSelf: "flex-start",
    boxShadow: "0 4px 15px rgba(250, 204, 21, 0.15)",
  },
};