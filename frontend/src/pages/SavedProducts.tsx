import React, { useState, useEffect } from "react";
import Header from "../components/general/Header";
import Footer from "../components/general/Footer";
import ResultsGrid, { Product } from "../components/mainpage/ResultsGrid";
import MainLayout from "../components/mainpage/MainLayout";

// SavedProducts get the saved products of the User
export default function SavedProducts() {
    const [savedProducts, setSavedProducts] = useState<Product[]>([]);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const fetchSavedProducts = async () => {
            try {
                // Make petition
                const response = await fetch("http://localhost:8080/api/user/saved-products", {
                    method: "GET",
                    credentials: "include",
                });

                if (response.ok) {
                    const data = await response.json();
                    setSavedProducts(data.data || []);
                } else {
                    console.error("Error al cargar los productos guardados");
                }
            } catch (error) {
                console.error("Error de conexión:", error);
            } finally {
                setIsLoading(false);
            }
        };

        fetchSavedProducts();
    }, []);

    return (
        <MainLayout
            header={<Header />}
            results={
                <>
                    <h2 style={{ fontSize: "28px", color: "#1f2937", marginBottom: "24px", marginTop: "24px", borderBottom: "2px solid #e5e7eb", paddingBottom: "10px" }}>
                        Productos Guardados:
                    </h2>

                    {/* Reuse of ResultsGrid */}
                    <ResultsGrid
                        products={savedProducts}
                        isLoading={isLoading}
                        emptyTitle="No tienes ningún producto guardado."
                        emptySubtitle="Explora y guarda productos para tenerlos siempre a mano."
                    />
                </>
            }
            footer={<Footer />}
        />
    );
}