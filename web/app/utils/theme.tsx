import { createContext, useContext, useEffect, useState, Dispatch, SetStateAction } from 'react';

type Theme = 'light' | 'dark';

interface ThemeContextType {
  theme: Theme;
  setTheme: Dispatch<SetStateAction<Theme>>;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

const getInitialTheme = (): Theme => {
  if (typeof window === 'undefined') return 'light'; // Default for SSR or non-browser env

  const storedPrefs = window.localStorage.getItem('color-theme');
  if (typeof storedPrefs === 'string' && (storedPrefs === 'light' || storedPrefs === 'dark')) {
    return storedPrefs as Theme;
  }

  const userMedia = window.matchMedia('(prefers-color-scheme: dark)');
  if (userMedia.matches) {
    return 'dark';
  }

  return 'light';
};

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setTheme] = useState<Theme>(getInitialTheme);

  useEffect(() => {
    const root = window.document.documentElement;
    const isDark = theme === "dark";

    root.classList.remove(isDark ? "light" : "dark");
    root.classList.add(theme);

    localStorage.setItem("color-theme", theme);
    // console.log(`Theme changed to: ${theme}`);
  }, [theme]);

  return (
    <ThemeContext.Provider value={{ theme, setTheme }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (context === undefined) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
} 