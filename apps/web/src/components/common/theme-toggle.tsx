import { Moon, Sun } from "lucide-react";
import { Button } from "@/components/ui/button";
import { useTheme } from "@/hooks/use-theme";

/**
 * Theme toggle button with sun/moon icons.
 */
export function ThemeToggle() {
	const { theme, toggleTheme } = useTheme();

	return (
		<Button
			variant="ghost"
			size="icon"
			onClick={toggleTheme}
			aria-label={`Switch to ${theme === "dark" ? "light" : "dark"} mode`}
			className="h-9 w-9"
		>
			{theme === "dark" ? (
				<Sun className="h-4 w-4" />
			) : (
				<Moon className="h-4 w-4" />
			)}
		</Button>
	);
}
