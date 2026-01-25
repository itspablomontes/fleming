import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { ArrowLeft } from "lucide-react";
import type { JSX } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { queryAuditEntries } from "@/features/audit/api";
import { AuditLogViewer } from "@/features/audit/components/audit-log-viewer";
import { AuditTargetType } from "@/features/audit/types";

import { ConsentDashboard } from "../components/consent-dashboard";
import { ConsentRequestWizard } from "../components/consent-request-wizard";

const auditPreviewLimit = 5;

export function ConsentPage(): JSX.Element {
	const auditQuery = useQuery({
		queryKey: ["audit-logs", "consent-preview"],
		queryFn: () =>
			queryAuditEntries({
				resourceType: AuditTargetType.Consent,
				limit: auditPreviewLimit,
			}),
	});

	return (
		<div className="min-h-screen bg-gray-50 px-4 py-6 dark:bg-gray-950 md:px-8 md:py-10">
			<div className="mx-auto flex w-full max-w-5xl flex-col gap-6">
				<div className="flex flex-col gap-3">
					<Button variant="ghost" size="sm" asChild className="w-fit">
						<Link to="/">
							<ArrowLeft className="h-4 w-4" />
							Back to timeline
						</Link>
					</Button>
					<div className="flex flex-col gap-2">
						<h1 className="text-2xl font-semibold text-gray-900 dark:text-gray-100">
							Consent
						</h1>
						<p className="text-sm text-muted-foreground">
							Control who can view or update your medical timeline.
						</p>
					</div>
				</div>

				<ConsentRequestWizard onSuccess={() => auditQuery.refetch()} />

				<ConsentDashboard />

				<Card className="border-border bg-white dark:bg-gray-900">
					<CardHeader className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
						<div className="space-y-1">
							<CardTitle className="text-base font-semibold text-foreground">
								Recent consent activity
							</CardTitle>
							<p className="text-sm text-muted-foreground">
								Review recent consent actions and audit trail entries.
							</p>
						</div>
						<Button variant="outline" asChild>
							<Link to="/audit">View full audit log</Link>
						</Button>
					</CardHeader>
					<CardContent>
						<AuditLogViewer
							entries={auditQuery.data ?? []}
							isLoading={auditQuery.isLoading}
							errorMessage={
								auditQuery.error ? "Failed to load audit entries." : undefined
							}
							onRefresh={() => auditQuery.refetch()}
						/>
					</CardContent>
				</Card>
			</div>
		</div>
	);
}
