import React from 'react'
import ReactDOM from 'react-dom/client'
import {
  createBrowserRouter,
  createRoutesFromElements,
  Route,
  RouterProvider,
} from "react-router-dom";
// import App from './App.tsx'
import './index.css'
import { Root } from './routes/root.tsx';
import ErrorPage from './error-page.js';
import { Login } from './credential-base/Login.tsx';
import { Dashboard } from './routes/Dashboard.tsx';

const router = createBrowserRouter(
  createRoutesFromElements(
    <Route path="/" element={<Root />} errorElement={<ErrorPage />}>
      <Route
        path="dashboard"
        element={<Dashboard />}
        loader={({ request }) =>
          fetch("/api/dashboard.json", {
            signal: request.signal,
          })
        }
      />
        <Route
          path="login"
          element={<Login />}
          // loader={redirectIfUser}
        />
        <Route path="logout" />
    </Route>
  ));

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
)
