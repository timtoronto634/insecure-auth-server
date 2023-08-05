import React from 'react'
import ReactDOM from 'react-dom/client'
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
// import App from './App.tsx'
import './index.css'
import { Root, loader as rootLoader } from './routes/root.tsx';
import ErrorPage from './error-page.tsx';
import Contact from './routes/contact.tsx';

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    errorElement:  <ErrorPage  />,
    loader: rootLoader,
    children: [
      {
        path: "contacts/:contactId",
        element: <Contact />,
      },
    ],
  }
]);

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
)
