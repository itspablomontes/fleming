import { useCallback, useEffect, useState } from "react";

type Theme = "dark" | "light";

const STORAGE_KEY = "fleming-theme";

/**
 * Theme hook with localStorage persistence.
 * Applies 'dark' class to documentElement for theme management.
 * Default is 'dark'.
 */
export function useTheme() {
	const [theme, setTheme] = useState<Theme>(() => {
		if (typeof window === "undefined") return "dark";
		const stored = localStorage.getItem(STORAGE_KEY);
		// Default to dark unless explicitly light
		return stored === "light" ? "light" : "dark";
	});

	useEffect(() => {
		const root = document.documentElement;
		if (theme === "dark") {
			root.classList.add("dark");
		} else {
			root.classList.remove("dark");
		}
		localStorage.setItem(STORAGE_KEY, theme);
	}, [theme]);

	const toggleTheme = useCallback(() => {
		setTheme((prev) => (prev === "dark" ? "light" : "dark"));
	}, []);

	return { theme, setTheme, toggleTheme };
}
