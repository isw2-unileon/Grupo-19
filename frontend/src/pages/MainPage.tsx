//import React from "react";
import { useState } from "react";
import MainLayout from "../components/mainpage/MainLayout";
import SearchBar from "../components/mainpage/SearchBar";
import ResultsGrid from "../components/mainpage/ResultsGrid";
import Header from "../components/general/Header";
import Footer from "../components/general/Footer";

export default function MainPage() {
  const [productsList, setProductsList] = useState([]);
  const [isLoading, setIsLoading] = useState(false);
  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

  // Method called by SearchBar when user press "Añadir" or Enter
  const procesarBusqueda = async (texto: string, esEnlace: boolean) => {
    setIsLoading(true);
    if (esEnlace) {
      console.log("Detectado un ENLACE. Llamando al scraper con:", texto);
      try {
        const response = await fetch(`${API_URL}/api/track`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
          body: JSON.stringify({ url: texto }),
        });

        const data = await response.json();

        if (response.ok) {
          setProductsList([data.data] as any);
        } else {
          console.error("Error al rastrear producto:", data.error);
        }
      } catch (error) {
        console.error("Error de red con el scraper:", error);
      }

    } else {
      console.log("Detectado TEXTO NORMAL. Buscando en base de datos:", texto);
      try {
        const response = await fetch(`${API_URL}/api/products/search?q=${encodeURIComponent(texto)}`, {
          method: "GET",
          credentials: "include",
        });

        const data = await response.json();

        if (response.ok) {
          // Refresh product list with the results
          setProductsList(data);
        } else {
          console.error("Error en la búsqueda:", data.error);
        }
      } catch (error) {
        console.error("Error de red en la búsqueda:", error);
      }
    }
    setIsLoading(false);
  };

  return (
    <MainLayout
      header={<Header />}
      search={<SearchBar onSearch={procesarBusqueda} />}
      results={<ResultsGrid products={productsList} isLoading={isLoading} />}
      footer={<Footer />}
    />
  );
}