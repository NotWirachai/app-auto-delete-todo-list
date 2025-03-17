import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import AppTodoList from "./AppTodoList.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <AppTodoList />
  </StrictMode>
);
