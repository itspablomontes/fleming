import { Loader2 } from "lucide-react";

/**
 * Full-screen loading state with skeleton animation.
 * Used during auth state resolution.
 */
export function LoadingScreen() {
    return (
        <div className="min-h-screen bg-background flex items-center justify-center">
            <div className="flex flex-col items-center gap-4">
                <Loader2 className="h-8 w-8 animate-spin text-primary" />
                <p className="text-muted-foreground text-sm">Loading...</p>
            </div>
        </div>
    );
}
