import React from "react";

type Props = {
  headerLeft: React.ReactNode;
  headerRight: React.ReactNode;
  search: React.ReactNode;
  results: React.ReactNode;
};

export default function MainLayout({ headerLeft, headerRight, search, results }: Props) {
  return (
    <div style={styles.container}>
      <header style={styles.headerContainer}>
        <div style={styles.headerLeft}>{headerLeft}</div>
        <div style={styles.headerRight}>{headerRight}</div>
      </header>
      <section style={styles.searchContainer}>{search}</section>
      <main style={styles.resultsContainer}>{results}</main>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    minHeight: "100vh",
    backgroundColor: "#fafafa",
    fontFamily: "sans-serif",
    padding: "40px 80px", 
    boxSizing: "border-box",
  },
  headerContainer: {
    display: "flex",
    justifyContent: "space-between",
    alignItems: "center",
    marginBottom: "60px",
    width: "100%",
  },
  headerLeft: {
    display: "flex",
    alignItems: "center",
  },
  headerRight: {
    display: "flex",
    gap: "20px",
    alignItems: "center",
  },
  searchContainer: {
    display: "flex",
    justifyContent: "center",
    marginBottom: "60px",
  },
  resultsContainer: {
    maxWidth: "1100px",
    margin: "0 auto",
    width: "100%",
  },
};