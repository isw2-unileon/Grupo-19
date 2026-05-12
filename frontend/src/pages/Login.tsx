import { useState } from "react";
import LoginLayout from "../components/login/LoginLayout";
import LoginForm from "../components/login/Loginform";
import LoginAnimation from "../components/login/LoginAnimation";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = (
    e: React.FormEvent,
    isLogin: boolean,
    name: string,
    confirmPassword: string
  ) => {
    e.preventDefault();

    if (isLogin) {
      // TODO: Cuando conectemos el backend, aquí irá el código de LOGIN
      console.log("Datos para LOGIN:", { email, password });
    } else {
      // Basic password validation
      if (password !== confirmPassword) {
        alert("Las contraseñas no coinciden");
        return;
      }

      // TODO: Cuando conectemos el backend, aquí irá el código de REGISTRO
      console.log("Datos para REGISTRO:", { name, email, password });
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