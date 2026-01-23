/**
 * Get a cookie value by name.
 */
export function getCookie(name: string): string | null {
	if (typeof document === "undefined") return null;
	const nameEQ = `${name}=`;
	const ca = document.cookie.split(";");
	for (let i = 0; i < ca.length; i++) {
		let c = ca[i];
		while (c.charAt(0) === " ") c = c.substring(1, c.length);
		if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
	}
	return null;
}

/**
 * Set a cookie value.
 */
export function setCookie(name: string, value: string, days = 7): void {
	if (typeof document === "undefined") return;
	let expires = "";
	if (days) {
		const date = new Date();
		date.setTime(date.getTime() + days * 24 * 60 * 60 * 1000);
		expires = `; expires=${date.toUTCString()}`;
	}
	// biome-ignore lint/suspicious/noDocumentCookie: Utility for managing cookies
	document.cookie = `${name}=${value || ""}${expires}; path=/; SameSite=Lax`;
}

/**
 * Delete a cookie by setting its expiration to the past.
 */
export function deleteCookie(name: string): void {
	if (typeof document === "undefined") return;
	// biome-ignore lint/suspicious/noDocumentCookie: Utility for managing cookies
	document.cookie = `${name}=; Max-Age=-99999999; path=/;`;
}
