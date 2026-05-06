import { useState } from "react";
import LoginLayout from "../components/login/LoginLayout";
import LoginForm from "../components/login/Loginform";
import LoginAnimation from "../components/login/LoginAnimation";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    console.log({ email, password });
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