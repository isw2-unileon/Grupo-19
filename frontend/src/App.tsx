import { useState, useEffect } from "react";
import { BrowserRouter, Routes, Route, Navigate, Outlet } from "react-router-dom";

import Login from "./pages/Login";
import MainPage from "./pages/MainPage";
import Profile from "./pages/Profile";
import Notifications from "./pages/Notifications";
import SavedProducts from "./pages/SavedProducts";
import Product from "./pages/Product";

function ProtectedRoute({ isAuth }: { isAuth: boolean }) {
  return isAuth ? <Outlet /> : <Navigate to="/" replace />;
}

export default function App() {
  const [isAuth, setIsAuth] = useState(false);
  const [isLoading, setIsLoading] = useState(true); // Estado de carga mientras comprueba el token

  useEffect(() => {
    async function checkSession() {
      try {
        // Llamamos al endpoint de perfil para ver si el backend nos reconoce
        const response = await fetch("http://localhost:8080/api/user/profile", {
          method: "GET",
          credentials: "include",
        });

        if (response.ok) {
          setIsAuth(true);
        } else {
          setIsAuth(false);
        }
      } catch (error) {
        console.error("Error al verificar la sesión:", error);
        setIsAuth(false);
      } finally {
        setIsLoading(false);
      }
    }

    checkSession();
  }, []);

  // Mientras comprueba la cookie, mostramos una pantalla limpia para que no parpadee el Login
  if (isLoading) {
    return (
      <div style={{ display: "flex", justifyContent: "center", alignItems: "center", height: "100vh", fontFamily: "sans-serif" }}>
        <p>Cargando aplicación...</p>
      </div>
    );
  }

  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/"
          element={
            !isAuth ? (
              <Login onLoginSuccess={() => setIsAuth(true)} />
            ) : (
              <Navigate to="/mainPage" replace />
            )
          }
        />

        <Route element={<ProtectedRoute isAuth={isAuth} />}>
          <Route path="/mainPage" element={<MainPage />} />
          <Route path="/profile" element={<Profile />} />
          <Route path="/notifications" element={<Notifications />} />
          <Route path="/savedProducts" element={<SavedProducts />} />
          <Route path="/product/:id" element={<Product />} />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}