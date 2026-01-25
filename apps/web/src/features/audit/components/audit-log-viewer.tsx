import type { JSX } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import type { AuditLogEntry } from "../types";
import { AuditLogEntryCard } from "./audit-log-entry";

interface AuditLogViewerProps {
	entries: AuditLogEntry[];
	isLoading?: boolean;
	errorMessage?: string;
	onRefresh?: () => void;
}

export function AuditLogViewer({
	entries,
	isLoading,
	errorMessage,
	onRefresh,
}: AuditLogViewerProps): JSX.Element {
	if (isLoading) {
		return (
			<Card className="bg-white dark:bg-gray-900">
				<CardHeader className="flex flex-row items-center justify-between">
					<CardTitle className="text-sm font-semibold text-gray-900 dark:text-gray-100">
						Audit Logs
					</CardTitle>
					{onRefresh && (
						<Button variant="outline" size="sm" onClick={onRefresh}>
							Refresh
						</Button>
					)}
				</CardHeader>
				<CardContent className="space-y-3">
					{Array.from({ length: 4 }).map((_, index) => (
						<Skeleton key={`audit-skeleton-loading-${Date.now()}-${index}`} className="h-20 w-full" />
					))}
				</CardContent>
			</Card>
		);
	}

	if (errorMessage) {
		return (
			<Card className="bg-white dark:bg-gray-900">
				<CardHeader>
					<CardTitle className="text-sm font-semibold text-gray-900 dark:text-gray-100">
						Audit Logs
					</CardTitle>
				</CardHeader>
				<CardContent className="text-sm text-destructive">
					{errorMessage}
				</CardContent>
			</Card>
		);
	}

	return (
		<Card className="bg-white dark:bg-gray-900">
			<CardHeader className="flex flex-row items-center justify-between">
				<CardTitle className="text-sm font-semibold text-gray-900 dark:text-gray-100">
					Audit Logs
				</CardTitle>
				{onRefresh && (
					<Button variant="outline" size="sm" onClick={onRefresh}>
						Refresh
					</Button>
				)}
			</CardHeader>
			<CardContent className="space-y-3">
				{entries.length === 0 ? (
					<p className="text-sm text-muted-foreground">No audit entries yet.</p>
				) : (
					entries.map((entry) => (
						<AuditLogEntryCard key={entry.id} entry={entry} />
					))
				)}
			</CardContent>
		</Card>
	);
}
