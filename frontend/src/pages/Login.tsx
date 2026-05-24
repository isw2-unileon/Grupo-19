import { useState } from "react";
import LoginLayout from "../components/login/LoginLayout";
import LoginForm from "../components/login/Loginform";
import LoginAnimation from "../components/login/LoginAnimation";

interface LoginProps {
  onLoginSuccess: () => void;
}

export default function Login({ onLoginSuccess }: LoginProps) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  // State for showing messages
  const [message, setMessage] = useState("");
  const [opacity, setOpacity] = useState(0);

  // Function to show messages
  const showMessage = (text: string) => {
    setMessage(text);
    setOpacity(1);
    setTimeout(() => setOpacity(0), 3500);
    setTimeout(() => setMessage(""), 4000);
  };

  const handleSubmit = async (
    e: React.FormEvent,
    isLogin: boolean,
    name: string,
    confirmPassword: string
  ) => {
    e.preventDefault();

    if (isLogin) {
      // Login API call
      try {
        const response = await fetch("http://localhost:8080/api/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ email, password }),
        });

        const data = await response.json();

        if (response.ok) {
          onLoginSuccess();
        } else {
          showMessage(data.error);
        }
      } catch (error) {
        showMessage("Error de conexión con el servidor");
      }
    } else {
      // Basic password validation
      if (password !== confirmPassword) {
        showMessage("Las contraseñas no coinciden");
        return;
      }

      // Register API call
      try {
        const response = await fetch("http://localhost:8080/api/register", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({ name, email, password }),
        });

        const data = await response.json();

        if (response.ok) {
          showMessage("¡Usuario registrado correctamente! Ya puedes iniciar sesión.");
        } else {
          showMessage(data.error);
        }
      } catch (error) {
        showMessage("Error de conexión con el servidor");
      }
    }
  };

  return (
    <LoginLayout
      left={
        <div
          style={{
            display: "flex",
            flexDirection: "column",
            alignItems: "center",
            position: "relative",
            width: "100%"
          }}
        >
          <LoginForm
            email={email}
            password={password}
            setEmail={setEmail}
            setPassword={setPassword}
            handleSubmit={handleSubmit}
          />

          {message && (
            <div
              style={{
                position: "absolute",
                top: "100%",
                marginTop: "20px",
                borderRadius: "12px",
                boxShadow: "0 10px 30px rgba(0,0,0,0.1)",
                padding: "20px",
                width: "320px",
                boxSizing: "border-box",
                backgroundColor: "rgba(255, 204, 0, 0.2)",
                border: "1px solid rgba(255, 204, 0, 0.5)",
                color: "#806600",
                opacity: opacity,
                transition: "opacity 0.5s ease-in-out",
                textAlign: "center",
                fontSize: "14px",
                fontWeight: "bold",
                pointerEvents: "none",
              }}
            >
              {message}
            </div>
          )}
        </div>
      }
      right={<LoginAnimation />}
    />
  );
}