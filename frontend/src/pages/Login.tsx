import { useState } from "react";
import LoginLayout from "../components/login/LoginLayout";
import LoginForm from "../components/login/Loginform";
import LoginAnimation from "../components/login/LoginAnimation";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

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
          alert("¡Login exitoso!");
        } else {
          // Muestra el mensaje de error que devuelve Gin (data.error)
          alert(data.error);
        }
      } catch (error) {
        alert("Error de conexión con el servidor");
      }
    } else {
      // Basic password validation
      if (password !== confirmPassword) {
        alert("Las contraseñas no coinciden");
        return;
      }

      // Register API call
      try {
        const response = await fetch("http://localhost:8080/api/register", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          // Enviamos "name", "email" y "password" para que coincida con el struct de Go
          body: JSON.stringify({ name, email, password }),
        });

        const data = await response.json();

        if (response.ok) {
          alert("¡Usuario registrado correctamente! Ya puedes iniciar sesión.");
        } else {
          alert(data.error);
        }
      } catch (error) {
        alert("Error de conexión con el servidor");
      }
    }
  };

  return (
    <LoginLayout
      left={
        <LoginForm
          email={email}
          password={password}
          setEmail={setEmail}
          setPassword={setPassword}
          handleSubmit={handleSubmit}
        />
      }
      right={<LoginAnimation />}
    />
  );
}