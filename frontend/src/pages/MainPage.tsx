import React, { useState } from "react";
import MainLayout from "../components/mainpage/MainLayout";
import SearchBar from "../components/mainpage/SearchBar";
import ResultsGrid from "../components/mainpage/ResultsGrid";
// Ajusta las rutas dependiendo de dónde guardaste los nuevos componentes
import Header from "../components/general/Header";
import Footer from "../components/general/Footer";

export default function MainPage() {
  const [productsList] = useState([]);

  return (
    <MainLayout
      header={<Header />}
      search={<SearchBar />}
      results={<ResultsGrid products={productsList} />}
      footer={<Footer />}
    />
  );
}