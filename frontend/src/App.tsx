import { useState } from "react";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";

import Login from "./pages/Login";
import MainPage from "./pages/MainPage";
import Profile from "./pages/Profile";
import Notifications from "./pages/Notifications";
import SavedProducts from "./pages/SavedProducts";

export default function App() {
  const [isAuth, setIsAuth] = useState(false);

  return (
    <BrowserRouter>
      {/* Lista de vistas posibles */}
      <Routes>

        {/* RUTA PRINCIPAL (/) */}
        {/* Si no está autenticado, muestra tu Login. Si ya lo está, lo desvía a /main */}
        <Route
          path="/"
          element={
            !isAuth ? (
              <Login onLoginSuccess={() => setIsAuth(true)} />
            ) : (
              <Navigate to="/mainPage" />
            )
          }
        />

        {/* RESTO DE RUTAS */}
        <Route path="/mainPage" element={<MainPage />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/notifications" element={<Notifications />} />
        <Route path="/savedProducts" element={<SavedProducts />} />

      </Routes>
    </BrowserRouter>
  );
}