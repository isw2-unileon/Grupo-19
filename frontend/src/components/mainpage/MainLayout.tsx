import React from "react";

type Props = {
  header: React.ReactNode;
  search?: React.ReactNode; // La '?' lo hace opcional por si una vista no tiene buscador
  results?: React.ReactNode; // También opcional
  footer: React.ReactNode;
};

export default function MainLayout({ header, search, results, footer }: Props) {
  return (
    <div style={styles.container}>
      {header}
      {/* Solo dibuja el buscador si la vista lo envía */}
      {search && <section style={styles.searchContainer}>{search}</section>}
      {/* Solo dibuja los resultados si la vista los envía */}
      {results && <main style={styles.resultsContainer}>{results}</main>}
      {footer}
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    minHeight: "100vh",
    backgroundColor: "#fafafa",
    fontFamily: "sans-serif",
    display: "flex",
    flexDirection: "column", // To organize the elements from top to bottom
  },
  searchContainer: {
    display: "flex",
    justifyContent: "center",
    marginBottom: "60px",
    padding: "40px 80px 0 80px", // Inside margins applied
  },
  resultsContainer: {
    maxWidth: "1100px",
    margin: "0 auto",
    width: "100%",
    padding: "0 80px 40px 80px",
    flex: 1, // Forces the central container to expand and push the footer towards the back
  },
};