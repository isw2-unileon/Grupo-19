import React, { useState, useEffect } from "react";
import Header from "../components/general/Header";
import Footer from "../components/general/Footer";

export default function Profile() {
  // --- For porfile fields ---
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [isLoadingProfile, setIsLoadingProfile] = useState(true);

  // --- States for password change ---
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  // --- Visual Feedback (Mensajes y Opacidades) ---
  const [profileMsg, setProfileMsg] = useState({ text: "", isError: false });
  const [passwordMsg, setPasswordMsg] = useState({ text: "", isError: false });

  const [profileOpacity, setProfileOpacity] = useState(0);
  const [passwordOpacity, setPasswordOpacity] = useState(0);

  // Funciones para mostrar mensajes que desaparecen (como en Login)
  const showProfileMessage = (text: string, isError: boolean) => {
    setProfileMsg({ text, isError });
    setProfileOpacity(1);
    setTimeout(() => setProfileOpacity(0), 3500);
    setTimeout(() => setProfileMsg({ text: "", isError: false }), 4000);
  };

  const showPasswordMessage = (text: string, isError: boolean) => {
    setPasswordMsg({ text, isError });
    setPasswordOpacity(1);
    setTimeout(() => setPasswordOpacity(0), 3500);
    setTimeout(() => setPasswordMsg({ text: "", isError: false }), 4000);
  };

  // LOAD PROFILE DATA FROM BACKEND (GET)
  useEffect(() => {
    async function fetchUserProfile() {
      try {
        const response = await fetch("http://localhost:8080/api/user/profile", {
          method: "GET",
          credentials: "include",
        });

        const data = await response.json();

        if (response.ok) {
          setUsername(data.Username);
          setEmail(data.Email);
        } else {
          showProfileMessage(data.error || "No se pudo cargar el perfil", true);
        }
      } catch (error) {
        console.error("Error al cargar perfil:", error);
        showProfileMessage("Error de conexión con el servidor", true);
      } finally {
        setIsLoadingProfile(false);
      }
    }

    fetchUserProfile();
  }, []);

  // UPDATE PROFLIE DATA (PUT)
  const handleUpdateProfile = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const response = await fetch("http://localhost:8080/api/user/profile", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ username, email }),
      });

      const data = await response.json();

      if (response.ok) {
        showProfileMessage("¡Perfil actualizado con éxito!", false);
      } else {
        showProfileMessage(data.error || "Error al actualizar", true);
      }
    } catch {
      showProfileMessage("Error de conexión con el servidor", true);
    }
  };

  // ACTUALIZAR CONTRASEÑA (PUT)
  const handleUpdatePassword = async (e: React.FormEvent) => {
    e.preventDefault();

    if (newPassword !== confirmPassword) {
      showPasswordMessage("Las contraseñas nuevas no coinciden", true);
      return;
    }

    try {
      const response = await fetch("http://localhost:8080/api/user/profile/password", {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({
          currentPassword,
          newPassword,
        }),
      });

      const data = await response.json();

      if (response.ok) {
        showPasswordMessage("¡Contraseña modificada con éxito!", false);
        setCurrentPassword("");
        setNewPassword("");
        setConfirmPassword("");
      } else {
        showPasswordMessage(data.error || "Error al cambiar contraseña", true);
      }
    } catch {
      showPasswordMessage("Error de conexión con el servidor", true);
    }
  };

  if (isLoadingProfile) {
    return (
      <div style={{ display: "flex", flexDirection: "column", minHeight: "100vh", backgroundColor: "#fafafa" }}>
        <Header />
        <main style={{ flex: 1, display: "flex", justifyContent: "center", alignItems: "center" }}>
          <p style={{ fontFamily: "sans-serif", fontWeight: "bold", color: "#4b5563" }}>Cargando datos de cuenta...</p>
        </main>
        <Footer />
      </div>
    );
  }

  return (
    <div style={{ display: "flex", flexDirection: "column", minHeight: "100vh", backgroundColor: "#fafafa", fontFamily: "sans-serif" }}>
      <Header />

      <main style={styles.mainContent}>
        <h2 style={styles.pageTitle}>Ajustes de Cuenta</h2>

        <div style={styles.sectionsGrid}>
          {/* PROFILE INFO */}
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

              {profileMsg.text && (
                <div style={{
                  ...styles.messageBox,
                  opacity: profileOpacity,
                  transition: "opacity 0.5s ease-in-out",
                  backgroundColor: profileMsg.isError ? "#fee2e2" : "rgba(255, 204, 0, 0.2)",
                  border: profileMsg.isError ? "1px solid #fca5a5" : "1px solid rgba(255, 204, 0, 0.5)",
                  color: profileMsg.isError ? "#991b1b" : "#806600"
                }}>
                  {profileMsg.text}
                </div>
              )}
            </form>
          </section>

          {/* PASSWORD AND SECURITY*/}
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

              {passwordMsg.text && (
                <div style={{
                  ...styles.messageBox,
                  opacity: passwordOpacity,
                  transition: "opacity 0.5s ease-in-out",
                  backgroundColor: passwordMsg.isError ? "#fee2e2" : "rgba(255, 204, 0, 0.2)",
                  border: passwordMsg.isError ? "1px solid #fca5a5" : "1px solid rgba(255, 204, 0, 0.5)",
                  color: passwordMsg.isError ? "#991b1b" : "#806600"
                }}>
                  {passwordMsg.text}
                </div>
              )}
            </form>
          </section>
        </div>
      </main>

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
    borderLeft: "5px solid #FACC15",
    paddingLeft: "12px",
    lineHeight: "1",
  },
  sectionsGrid: {
    display: "grid",
    gridTemplateColumns: "repeat(auto-fit, minmax(450px, 1fr))",
    gap: "32px",
  },
  card: {
    background: "white",
    borderRadius: "12px",
    padding: "32px",
    boxShadow: "0 4px 15px rgba(0,0,0,0.01)",
    border: "1px solid #eee",
    borderTop: "4px solid #FACC15",
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
  messageBox: {
    padding: "12px",
    borderRadius: "8px",
    fontSize: "13px",
    fontWeight: "bold",
    textAlign: "center",
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