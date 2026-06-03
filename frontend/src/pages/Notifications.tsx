import React, { useState, useEffect } from "react";
import Header from "../components/general/Header";
import Footer from "../components/general/Footer";

interface NotificationItem {
  id: number;
  type: string; 
  title: string;
  description: string;
  time: string;
  isRead: boolean;
}

export default function Notifications() {
  const [notifications, setNotifications] = useState<NotificationItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [errorMsg, setErrorMsg] = useState("");


  // LOAD NOTIFICATIONS (GET)
  useEffect(() => {
    async function fetchNotifications() {
      try {
        const response = await fetch("http://localhost:8080/api/user/notifications", {
          method: "GET",
          credentials: "include", 
        });

        const data = await response.json();

        if (response.ok) {
          setNotifications(data || []); 
        } else {
          setErrorMsg(data.error || "Error al cargar las notificaciones");
        }
      } catch (error) {
        console.error("Error cross-origin o de red:", error);
        setErrorMsg("Error de conexión con el servidor");
      } finally {
        setIsLoading(false);
      }
    }

    fetchNotifications();
  }, []);


  // LOGIC: MARK AS READ (PATCH)
  const markAsRead = async (id: number) => {
    try {
      const response = await fetch(`http://localhost:8080/api/user/notifications/${id}`, {
        method: "PATCH",
        credentials: "include", // Envía el token JWT en las cookies
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ isRead: true }),
      });

      if (response.ok) {
        setNotifications(notifications.map((n) => (n.id === id ? { ...n, isRead: true } : n)));
      } else {
        const data = await response.json();
        alert(data.error || "No se pudo marcar la notificación como leída");
      }
    } catch (error) {
      console.error("Error al actualizar la notificación:", error);
      alert("Error de conexión al marcar como leída");
    }
  };


  // DELETE NOTIFICATION (DELETE)
  const deleteNotification = async (id: number) => {
    // Pedimos confirmación al usuario por seguridad
    if (!confirm("¿Estás seguro de que deseas eliminar esta notificación?")) {
      return;
    }

    try {
      const response = await fetch(`http://localhost:8080/api/user/notifications/${id}`, {
        method: "DELETE",
        credentials: "include", // Envía el token JWT en las cookies
      });

      if (response.ok) {
        setNotifications(notifications.filter((n) => n.id !== id));
      } else {
        const data = await response.json();
        alert(data.error || "No se pudo eliminar la notificación");
      }
    } catch (error) {
      console.error("Error al eliminar la notificación:", error);
      alert("Error de conexión al eliminar la notificación");
    }
  };

  // ==========================================
  // CONFIGURACIÓN DINÁMICA POR TIPO DE ALERTA
  // ==========================================
  const getTypeStyles = (type: string) => {
    switch (type) {
      case "target_price":
        return { 
          color: "#ffffff", 
          bg: "#dc2626", 
          label: "🎯 PRECIO OBJETIVO ALCANZADO",
          borderColor: "#dc2626"
        };
      case "price_drop":
        return { 
          color: "#166534", 
          bg: "#dcfce7", 
          label: "📉 BAJADA DE PRECIO",
          borderColor: "#22c55e"
        };
      default:
        return { 
          color: "#854d0e", 
          bg: "#fef9c3", 
          label: "⚙️ SISTEMA",
          borderColor: "#eab308" 
        };
    }
  };

  if (isLoading) {
    return (
      <div style={{ display: "flex", flexDirection: "column", minHeight: "100vh", backgroundColor: "#fafafa" }}>
        <Header />
        <main style={{ flex: 1, display: "flex", justifyContent: "center", alignItems: "center", fontFamily: "sans-serif" }}>
          <p>Cargando tus alertas...</p>
        </main>
        <Footer />
      </div>
    );
  }

  return (
    <div style={{ display: "flex", flexDirection: "column", minHeight: "100vh", backgroundColor: "#fafafa", fontFamily: "sans-serif" }}>
      <Header />
      <main style={styles.mainContent}>
        <div style={styles.pageHeader}>
          <h2 style={styles.pageTitle}>Centro de Notificaciones</h2>
        </div>

        {errorMsg ? (
          <div style={{ color: "#991b1b", backgroundColor: "#fee2e2", padding: "16px", borderRadius: "8px", textAlign: "center", fontWeight: "bold" }}>
            {errorMsg}
          </div>
        ) : (
          <div style={styles.container}>
            {notifications.length === 0 ? (
              <div style={styles.emptyState}>
                <span style={{ fontSize: "48px" }}>🔔</span>
                <p style={styles.emptyText}>No tienes ninguna notificación por el momento.</p>
                <p style={{ color: "#6b7280", fontSize: "14px" }}>Te avisaremos aquí cuando los precios de tus productos favoritos cambien.</p>
              </div>
            ) : (
              <div style={styles.list}>
                {notifications.map((notification) => {
                  const badge = getTypeStyles(notification.type);
                  
                  const getCardBg = () => {
                    if (notification.isRead) return "#ffffff";
                    return notification.type === "target_price" ? "#fef2f2" : "#fffbeb";
                  };

                  return (
                    <div
                      key={notification.id}
                      style={{
                        ...styles.notificationCard,
                        borderLeft: notification.isRead ? "4px solid #d1d5db" : `4px solid ${badge.borderColor}`,
                        backgroundColor: getCardBg(),
                      }}
                    >
                      <div style={styles.cardBody}>
                        <div style={styles.badgeRow}>
                          <span style={{ ...styles.badge, backgroundColor: badge.bg, color: badge.color }}>
                            {badge.label}
                          </span>
                          <span style={styles.timeText}>
                            {new Date(notification.time).toLocaleDateString()} a las{" "}
                            {new Date(notification.time).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                          </span>
                        </div>
                        <h4 style={{ ...styles.notificationTitle, fontWeight: notification.isRead ? "600" : "bold" }}>
                          {notification.title}
                        </h4>
                        <p style={styles.notificationDescription}>{notification.description}</p>
                      </div>

                      <div style={styles.cardActions}>
                        {!notification.isRead && (
                          <button 
                            onClick={() => markAsRead(notification.id)} 
                            style={styles.actionButton}
                            title="Marcar como leído"
                          >
                            👁️
                          </button>
                        )}
                        <button 
                          onClick={() => deleteNotification(notification.id)} 
                          style={{ ...styles.actionButton, color: "#ef4444" }}
                          title="Eliminar"
                        >
                          🗑️
                        </button>
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        )}
      </main>
      <Footer />
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  mainContent: {
    flex: 1,
    width: "100%",
    maxWidth: "900px",
    margin: "0 auto",
    padding: "40px 20px 80px 20px",
    boxSizing: "border-box",
  },
  pageHeader: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    marginBottom: "32px",
  },
  pageTitle: {
    fontSize: "26px",
    color: "#1f2937",
    margin: 0,
    fontWeight: "bold",
    borderLeft: "5px solid #FACC15",
    paddingLeft: "12px",
    lineHeight: "1",
  },
  container: {
    width: "100%",
  },
  list: {
    display: "flex",
    flexDirection: "column",
    gap: "16px",
  },
  notificationCard: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    padding: "20px 24px",
    borderRadius: "8px",
    boxShadow: "0 2px 8px rgba(0,0,0,0.02)",
    border: "1px solid #eee",
    transition: "all 0.2s ease-in-out",
  },
  cardBody: {
    flex: 1,
    paddingRight: "20px",
  },
  badgeRow: {
    display: "flex",
    alignItems: "center",
    gap: "12px",
    marginBottom: "8px",
  },
  badge: {
    fontSize: "11px",
    fontWeight: "bold",
    padding: "4px 10px",
    borderRadius: "4px",
    letterSpacing: "0.5px",
  },
  timeText: {
    fontSize: "12px",
    color: "#9ca3af",
  },
  notificationTitle: {
    margin: "0 0 6px 0",
    fontSize: "16px",
    color: "#1f2937",
  },
  notificationDescription: {
    margin: 0,
    fontSize: "14px",
    color: "#4b5563",
    lineHeight: "1.4",
  },
  cardActions: {
    display: "flex",
    gap: "8px",
  },
  actionButton: {
    background: "#f3f4f6",
    border: "none",
    borderRadius: "6px",
    width: "36px",
    height: "36px",
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    cursor: "pointer",
    fontSize: "14px",
    transition: "background 0.2s",
  },
  emptyState: {
    textAlign: "center",
    padding: "80px 20px",
    background: "white",
    borderRadius: "12px",
    border: "1px solid #eee",
    boxShadow: "0 4px 15px rgba(0,0,0,0.01)",
  },
  emptyText: {
    fontSize: "18px",
    fontWeight: "bold",
    color: "#374151",
    marginTop: "16px",
    marginBottom: "8px",
  },
};