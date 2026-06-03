import React, { useState } from "react";

type Props = {
  onSearch: (texto: string, esEnlace: boolean) => void;
};

export default function SearchBar({ onSearch }: Props) {

  const [query, setQuery] = useState("");

  const handleAction = () => {
    if (!query.trim()) return;

    // Rule to know if is an URL or a name
    const urlPattern = /^(https?:\/\/|www\.)/i;
    const esEnlace = urlPattern.test(query.trim());

    onSearch(query.trim(), esEnlace);

    // Clen SearchBar
    setQuery("");
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      handleAction();
    }
  };

  return (
    <div style={styles.container}>
      {/* --- Sear bar --- */}
      <div style={styles.searchWrapper}>
        <div style={styles.iconContainer}>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={2.5}
            stroke="#9ca3af"
            style={styles.magnifier}
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              pathLength="1"
              d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
            />
          </svg>
        </div>

        {/* Text input */}
        <input
          type="text"
          placeholder="Busca un producto o introduce un enlace..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={handleKeyDown}
          style={styles.input}
        />
      </div>

      {/* --- Add button --- */}
      <button type="button" onClick={handleAction} style={styles.addButton}>
        Añadir
      </button>
    </div>
  );
}

const styles: Record<string, React.CSSProperties> = {
  container: {
    display: "flex",
    width: "100%",
    maxWidth: "750px",
    gap: "16px",
    alignItems: "center",
  },
  searchWrapper: {
    display: "flex",
    flex: 1,
    boxShadow: "0 4px 20px rgba(0,0,0,0.05)",
    borderRadius: "8px",
    border: "1px solid #ddd",
    backgroundColor: "white",
    overflow: "hidden",
    alignItems: "center",
  },
  iconContainer: {
    paddingLeft: "16px",
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
  },
  magnifier: {
    width: "20px",
    height: "20px",
  },
  input: {
    flex: 1,
    padding: "16px 16px 16px 12px",
    fontSize: "15px",
    border: "none",
    backgroundColor: "transparent",
    outline: "none",
  },
  addButton: {
    padding: "16px 32px",
    border: "none",
    background: "#FACC15",
    color: "black",
    fontWeight: "bold",
    fontSize: "15px",
    borderRadius: "8px",
    cursor: "pointer",
    boxShadow: "0 4px 15px rgba(250, 204, 21, 0.2)",
  },
};