import { useCallback, useEffect, useState } from "react";

type Theme = "dark" | "light";

const STORAGE_KEY = "fleming-theme";

/**
 * Theme hook with localStorage persistence.
 * Applies theme class to documentElement.
 */
export function useTheme() {
	const [theme, setTheme] = useState<Theme>(() => {
		if (typeof window === "undefined") return "dark";
		const stored = localStorage.getItem(STORAGE_KEY);
		return stored === "light" ? "light" : "dark";
	});

	useEffect(() => {
		const root = document.documentElement;
		root.classList.remove("dark", "light");
		root.classList.add(theme);
		localStorage.setItem(STORAGE_KEY, theme);
	}, [theme]);

	const toggleTheme = useCallback(() => {
		setTheme((prev) => (prev === "dark" ? "light" : "dark"));
	}, []);

	return { theme, setTheme, toggleTheme };
}
