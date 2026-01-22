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
 * Delete a cookie by setting its expiration to the past.
 */
export function deleteCookie(name: string): void {
    if (typeof document === "undefined") return;
    document.cookie = `${name}=; Max-Age=-99999999; path=/;`;
}
