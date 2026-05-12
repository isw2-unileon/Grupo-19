import { useState } from "react";

type Props = {
  email: string;
  password: string;
  setEmail: (v: string) => void;
  setPassword: (v: string) => void;
  handleSubmit: (e: React.FormEvent, isLogin: boolean, name: string, confirmPassword: string) => void;
};

export default function LoginForm({
  email,
  password,
  setEmail,
  setPassword,
  handleSubmit,
}: Props) {
  // State to control the visible form
  const [isLogin, setIsLogin] = useState(true);

  const [name, setName] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  return (
    <form onSubmit={(e) => handleSubmit(e, isLogin, name, confirmPassword)} style={styles.form}>
      <h2 style={{ marginBottom: "20px", textAlign: "center", fontWeight: "bold" }}>
        {isLogin ? "Inicia sesión en tu cuenta" : "Crea una cuenta"}
      </h2>

      {/* Only for register */}
      {!isLogin && (
        <input
          type="text"
          placeholder="Nombre"
          value={name}
          onChange={(e) => setName(e.target.value)}
          style={styles.input}
        />
      )}

      <input
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        style={styles.input}
      />

      <input
        type="password"
        placeholder="Contraseña"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        style={styles.input}
      />

      {/* Only for register */}
      {!isLogin && (
        <input
          type="password"
          placeholder="Confirmar Contraseña"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          style={styles.input}
        />
      )}

      {/* Button to submit form */}
      <button type="submit" style={styles.button}>
        {isLogin ? "Entrar" : "Registrarse"}
      </button>

      {/* Button to change form */}
      <div style={styles.switchContainer}>
        <span style={{ fontSize: "14px" }}>
          {isLogin ? "¿No tienes cuenta? " : "¿Ya tienes cuenta? "}
        </span>
        <button
          type="button"
          onClick={() => {
            setIsLogin(!isLogin);
            setEmail("");
            setPassword("");
            setName("");
            setConfirmPassword("");
          }}
          style={styles.linkButton}
        >
          {isLogin ? "Regístrate aquí" : "Inicia sesión"}
        </button>
      </div>
    </form>
  );
}

// Styles
const styles: Record<string, React.CSSProperties> = {
  form: {
    display: "flex",
    flexDirection: "column",
    gap: "12px",
    width: "320px",
    padding: "30px",
    background: "white",
    borderRadius: "12px",
    boxShadow: "0 10px 30px rgba(0,0,0,0.1)",
  },
  input: {
    padding: "12px",
    borderRadius: "6px",
    border: "1px solid #ddd",
    fontSize: "14px",
  },
  button: {
    padding: "12px",
    borderRadius: "6px",
    border: "none",
    background: "#FACC15",
    color: "black",
    cursor: "pointer",
    fontWeight: "bold",
  },
  switchContainer: {
    textAlign: "center",
    marginTop: "10px",
  },
  linkButton: {
    background: "none",
    border: "none",
    color: "#000",
    textDecoration: "underline",
    cursor: "pointer",
    fontWeight: "bold",
    fontSize: "14px",
    padding: 0,
  },
};