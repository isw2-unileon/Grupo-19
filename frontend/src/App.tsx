import { useState } from "react";
import Login from "./pages/Login";
import Home from "./pages/Home";

export default function App() {
  // State for control the login status
  const [isAuth, setIsAuth] = useState(false);

  // If isAuth is true we redirect to the home page
  if (isAuth) {
    return <Home />;
  }

  // When login page detects the login process is finished, it set isAuth true and React redirect to the home page.
  return <Login onLoginSuccess={() => setIsAuth(true)} />;
}