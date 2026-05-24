import React, { useState } from "react";
import MainLayout from "../components/mainpage/MainLayout";
import SearchBar from "../components/mainpage/SearchBar";
import ResultsGrid from "../components/mainpage/ResultsGrid";
import HeaderActions from "../components/mainpage/HeaderActions";
import logoImage from "../assets/Logo.png"; 

export default function Home() {
  // Usamos solo productsList para que ESLint no se queje de variables sin usar
  const [productsList] = useState([]);

  return (
    <MainLayout
      headerLeft={
        <div style={styles.logoContainer}>
          <img 
            src={logoImage} 
            alt="ProTracker Logo" 
            style={styles.logoImg} 
          />
        </div>
      }
      headerRight={<HeaderActions />}
      search={<SearchBar />}
      /* Pasamos el array vacío de forma segura */
      results={<ResultsGrid products={productsList} />} 
    />
  );
}

const styles: Record<string, React.CSSProperties> = {
  logoContainer: {
    display: "flex",
    alignItems: "center",
    justifyContent: "center",
    padding: "5px 0",
  },
  logoImg: {
    height: "100px",       
    width: "auto",        
    objectFit: "contain",
  },
};