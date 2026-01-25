import type { JSX } from "react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import type { AuditLogEntry } from "../types";

interface AuditLogEntryProps {
	entry: AuditLogEntry;
}

const formatTimestamp = (timestamp: Date) => timestamp.toLocaleString();

export function AuditLogEntryCard({ entry }: AuditLogEntryProps): JSX.Element {
	return (
		<Card className="bg-white dark:bg-gray-900">
			<CardHeader className="flex flex-row items-start justify-between gap-4">
				<div className="space-y-1">
					<CardTitle className="text-sm font-semibold text-gray-900 dark:text-gray-100">
						{entry.action}
					</CardTitle>
					<p className="text-xs text-muted-foreground">
						{entry.resourceType} • {entry.resourceId}
					</p>
				</div>
				<Badge variant="secondary" className="text-xs">
					{formatTimestamp(entry.timestamp)}
				</Badge>
			</CardHeader>
			<CardContent className="space-y-2 text-xs text-muted-foreground">
				<p>
					<span className="font-medium text-gray-700 dark:text-gray-200">
						Actor:
					</span>{" "}
					{entry.actor}
				</p>
				{entry.hash && (
					<p className="truncate">
						<span className="font-medium text-gray-700 dark:text-gray-200">
							Hash:
						</span>{" "}
						{entry.hash}
					</p>
				)}
				{entry.previousHash && (
					<p className="truncate">
						<span className="font-medium text-gray-700 dark:text-gray-200">
							Prev:
						</span>{" "}
						{entry.previousHash}
					</p>
				)}
			</CardContent>
		</Card>
	);
}
