import ReactDOM from "react-dom/client";
import { RouterProvider } from "react-router-dom";

import { QueryProvider } from "./app/providers/query-provider";
import { router } from "./app/router/router";

import "./index.css";

ReactDOM.createRoot(
  document.getElementById("root")!
).render(
  <QueryProvider>
    <RouterProvider router={router} />
  </QueryProvider>
);