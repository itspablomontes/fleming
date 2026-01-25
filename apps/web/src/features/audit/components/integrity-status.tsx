import type { JSX } from "react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

interface IntegrityStatusProps {
	isValid?: boolean;
	message?: string;
}

export function IntegrityStatus({
	isValid,
	message,
}: IntegrityStatusProps): JSX.Element {
	const statusLabel =
		isValid === undefined ? "Not verified" : isValid ? "Valid" : "Invalid";

	return (
		<Card className="bg-white dark:bg-gray-900">
			<CardHeader className="flex flex-row items-start justify-between gap-4">
				<CardTitle className="text-sm font-semibold text-gray-900 dark:text-gray-100">
					Audit Chain Integrity
				</CardTitle>
				<Badge variant={isValid ? "default" : "secondary"}>{statusLabel}</Badge>
			</CardHeader>
			<CardContent className="text-xs text-muted-foreground">
				{message ?? "Run verification to confirm the chain integrity."}
			</CardContent>
		</Card>
	);
}
