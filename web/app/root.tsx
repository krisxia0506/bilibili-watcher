import {
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
} from "@remix-run/react";
import type { LinksFunction } from "@remix-run/node";
import { ThemeProvider, useTheme } from "~/utils/theme";

import "./tailwind.css";

export const links: LinksFunction = () => [
  { rel: "preconnect", href: "https://fonts.googleapis.com" },
  {
    rel: "preconnect",
    href: "https://fonts.gstatic.com",
    crossOrigin: "anonymous",
  },
  {
    rel: "stylesheet",
    href: "https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&display=swap",
  },
];

function InnerLayout({ children }: { children: React.ReactNode }) {
  const { theme, setTheme } = useTheme();

  const toggleTheme = () => {
    setTheme((prevTheme) => (prevTheme === "dark" ? "light" : "dark"));
  };

  // Determine theme colors for meta tag
  // These should match your actual body background colors for light and dark modes
  const lightThemeColor = "#ffffff"; // Assuming body bg-white in light mode
  const darkThemeColor = "#111827";  // Assuming body dark:bg-gray-900 in dark mode

  return (
    <html lang="en" className={typeof window !== 'undefined' ? document.documentElement.className : 'light'}>
      <head>
        <meta charSet="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        {/* Dynamically set theme-color based on the current theme */}
        <meta name="theme-color" content={theme === 'dark' ? darkThemeColor : lightThemeColor} />
        <Meta />
        <Links />
      </head>
      <body className="bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 transition-colors duration-200">
        <div className="p-4 flex justify-end">
          <button 
            onClick={toggleTheme} 
            className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200 rounded-md hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors duration-200"
          >
            Switch to {theme === "light" ? "Dark" : "Light"} Mode
          </button>
        </div>
        {children}
        <ScrollRestoration />
        <Scripts />
      </body>
    </html>
  );
}

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <ThemeProvider>
      <InnerLayout>{children}</InnerLayout>
    </ThemeProvider>
  );
}

export default function App() {
  return <Outlet />;
}
