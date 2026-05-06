type Props = {
  email: string;
  password: string;
  setEmail: (v: string) => void;
  setPassword: (v: string) => void;
  handleSubmit: (e: React.FormEvent) => void;
};

export default function LoginForm({
  email,
  password,
  setEmail,
  setPassword,
  handleSubmit,
}: Props) {
  return (
    <form onSubmit={handleSubmit} style={styles.form}>
      <h2 style={{ marginBottom: "20px" }}>
        Inicia sesión en tu cuenta
      </h2>

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

      <button type="submit" style={styles.button}>
        Entrar
      </button>
    </form>
  );
}

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
};