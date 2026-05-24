import { useState } from "react";
import Login from "./pages/Login";
import MainPage from "./pages/MainPage";

export default function App() {
  // State for control the login status
  const [isAuth, setIsAuth] = useState(false);

  // If isAuth is true we redirect to the home page
  if (isAuth) {
    return <MainPage />;
  }

  // When login page detects the login process is finished, it set isAuth true and React redirect to the home page.
  return <Login onLoginSuccess={() => setIsAuth(true)} />;
}