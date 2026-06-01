export default function LoginLayout({
  left,
  right,
}: {
  left: React.ReactNode;
  right: React.ReactNode;
}) {
  return (
    <div style={styles.container}>
      <div style={styles.left}>{left}</div>
      <div style={styles.right}>{right}</div>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    display: "flex",
    height: "100vh",
    fontFamily: "sans-serif",
  },

  left: {
    flex: 4,
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    background: "#ffffff",
  },

  right: {
    flex: 6,
    background: "linear-gradient(135deg, #FACC15, #f59e0b)",
    position: "relative",
    overflow: "hidden",
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
  },
};