import { useState, useEffect } from "react";
import { BrowserRouter, Routes, Route, Navigate, Outlet } from "react-router-dom";

import Login from "./pages/Login";
import MainPage from "./pages/MainPage";
import Profile from "./pages/Profile";
import Notifications from "./pages/Notifications";
import SavedProducts from "./pages/SavedProducts";
import Product from "./pages/Product";
import About from "./pages/About";
import Contact from "./pages/Contact";
import Terms from "./pages/Terms";
import Privacy from "./pages/Privacy";

function ProtectedRoute({ isAuth }: { isAuth: boolean }) {
  return isAuth ? <Outlet /> : <Navigate to="/" replace />;
}

export default function App() {
  const [isAuth, setIsAuth] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

  useEffect(() => {
    async function checkSession() {
      try {
        const response = await fetch(`${API_URL}/api/user/profile`, {
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

  // Clean interface while cookie is checked
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
          <Route path="/about" element={<About />} />
          <Route path="/contact" element={<Contact />} />
          <Route path="/terms" element={<Terms />} />
          <Route path="/privacy" element={<Privacy />} />
        </Route>

        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}