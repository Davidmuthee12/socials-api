import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import App, { ActivateAccountPage } from "./App.tsx";

export function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<App />} />
        <Route
          path="/users/activate/:token"
          element={<ActivateAccountPage />}
        />
        <Route path="/users/activate" element={<Navigate to="/" replace />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  );
}
